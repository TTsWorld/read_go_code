// Copyright (c) 2023 Uber Technologies, Inc.
// 版权归Uber Technologies公司所有
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// 特此免费授予任何获得本软件副本的人使用、复制、修改、分发等权利
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// 上述版权声明和许可声明应包含在软件的所有副本中
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
// 本软件"按原样"提供，不提供任何明示或暗示的保证

package ztest // ztest包：测试工具

import (
	"sort" // sort包：排序工具
	"sync" // sync包：同步原语
	"time" // time包：时间处理
)

// MockClock is a fake source of time.
// It implements standard time operations,
// but allows the user to control the passage of time.
//
// Use the [Add] method to progress time.
// MockClock是一个虚假的时间源。它实现标准的时间操作，
// 但允许用户控制时间的流逝。
//
// 使用[Add]方法来推进时间。
type MockClock struct {
	mu  sync.RWMutex // 读写互斥锁
	now time.Time    // 当前时间

	// The MockClock works by maintaining a list of waiters.
	// Each waiter knows the time at which it should be resolved.
	// When the clock advances, all waiters that are in range are resolved
	// in chronological order.
	// MockClock通过维护一个等待者列表来工作。
	// 每个等待者都知道它应该被解决的时间。
	// 当时钟前进时，所有在范围内的等待者都会按时间顺序被解决。
	waiters []waiter // 等待者列表
}

// NewMockClock builds a new mock clock
// using the current actual time as the initial time.
// NewMockClock使用当前实际时间作为初始时间构建新的模拟时钟。
func NewMockClock() *MockClock {
	return &MockClock{
		now: time.Now(), // 设置当前时间
	}
}

// Now reports the current time.
// Now报告当前时间。
func (c *MockClock) Now() time.Time {
	c.mu.RLock()         // 获取读锁
	defer c.mu.RUnlock() // 延迟释放读锁
	return c.now         // 返回当前时间
}

// NewTicker returns a time.Ticker that ticks at the specified frequency.
//
// As with [time.NewTicker],
// the ticker will drop ticks if the receiver is slow,
// and the channel is never closed.
//
// Calling Stop on the returned ticker is a no-op.
// The ticker only runs when the clock is advanced.
func (c *MockClock) NewTicker(d time.Duration) *time.Ticker {
	ch := make(chan time.Time, 1)

	var tick func(time.Time)
	tick = func(now time.Time) {
		next := now.Add(d)
		c.runAt(next, func() {
			defer tick(next)

			select {
			case ch <- next:
				// ok
			default:
				// The receiver is slow.
				// Drop the tick and continue.
			}
		})
	}
	tick(c.Now())

	return &time.Ticker{C: ch}
}

// runAt schedules the given function to be run at the given time.
// The function runs without a lock held, so it may schedule more work.
func (c *MockClock) runAt(t time.Time, fn func()) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.waiters = append(c.waiters, waiter{until: t, fn: fn})
}

type waiter struct { // 等待者结构体
	until time.Time // 等待直到的时间
	fn    func()    // 执行的函数
}

// Add progresses time by the given duration.
// Other operations waiting for the time to advance
// will be resolved if they are within range.
//
// Side effects of operations waiting for the time to advance
// will take effect on a best-effort basis.
// Avoid racing with operations that have side effects.
//
// Panics if the duration is negative.
// Add按给定的持续时间推进时间。等待时间推进的其他操作
// 如果在范围内将被解决。
//
// 等待时间推进的操作的副作用将在尽力而为的基础上生效。
// 避免与有副作用的操作竞争。
//
// 如果持续时间为负数，则panic。
func (c *MockClock) Add(d time.Duration) {
	if d < 0 { // 检查持续时间是否为负
		panic("cannot add negative duration") // 不能添加负数持续时间
	}

	c.mu.Lock()         // 获取写锁
	defer c.mu.Unlock() // 延迟释放写锁

	sort.Slice(c.waiters, func(i, j int) bool { // 按时间排序等待者
		return c.waiters[i].until.Before(c.waiters[j].until)
	})

	newTime := c.now.Add(d) // 计算新时间
	// newTime won't be recorded until the end of this method.
	// This ensures that any waiters that are resolved
	// are resolved at the time they were expecting.
	// newTime直到此方法结束才会被记录。
	// 这确保任何被解决的等待者都在它们期望的时间被解决。

	for len(c.waiters) > 0 { // 处理所有等待者
		w := c.waiters[0]           // 获取第一个等待者
		if w.until.After(newTime) { // 如果等待时间在新时间之后
			break // 跳出循环
		}
		c.waiters[0] = waiter{}   // 避免内存泄漏
		c.waiters = c.waiters[1:] // 移除第一个等待者

		// The waiter is within range.
		// Travel to the time of the waiter and resolve it.
		// 等待者在范围内。移动到等待者的时间并解决它。
		c.now = w.until // 设置当前时间为等待者的时间

		// The waiter may schedule more work
		// so we must release the lock.
		// 等待者可能安排更多工作，所以我们必须释放锁。
		c.mu.Unlock() // 释放锁
		w.fn()        // 执行等待者函数
		// Sleeping here is necessary to let the side effects of waiters
		// take effect before we continue.
		// 在这里睡眠是必要的，让等待者的副作用在我们继续之前生效。
		time.Sleep(1 * time.Millisecond) // 睡眀1毫秒
		c.mu.Lock()                      // 重新获取锁
	}

	c.now = newTime // 设置最终时间
}
