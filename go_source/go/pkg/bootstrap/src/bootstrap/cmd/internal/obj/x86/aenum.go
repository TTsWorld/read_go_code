// Code generated by go tool dist; DO NOT EDIT.
// This is a bootstrap copy of /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/internal/obj/x86/aenum.go

//line /Users/litiantian/code/m_code/read_go_code/go_source/go/src/cmd/internal/obj/x86/aenum.go:1
// Code generated by x86avxgen. DO NOT EDIT.

package x86

import "bootstrap/cmd/internal/obj"

//go:generate go run ../stringer.go -i $GOFILE -o anames.go -p x86

const (
	AAAA = obj.ABaseAMD64 + obj.A_ARCHSPECIFIC + iota
	AAAD
	AAAM
	AAAS
	AADCB
	AADCL
	AADCQ
	AADCW
	AADCXL
	AADCXQ
	AADDB
	AADDL
	AADDPD
	AADDPS
	AADDQ
	AADDSD
	AADDSS
	AADDSUBPD
	AADDSUBPS
	AADDW
	AADJSP
	AADOXL
	AADOXQ
	AAESDEC
	AAESDECLAST
	AAESENC
	AAESENCLAST
	AAESIMC
	AAESKEYGENASSIST
	AANDB
	AANDL
	AANDNL
	AANDNPD
	AANDNPS
	AANDNQ
	AANDPD
	AANDPS
	AANDQ
	AANDW
	AARPL
	ABEXTRL
	ABEXTRQ
	ABLENDPD
	ABLENDPS
	ABLENDVPD
	ABLENDVPS
	ABLSIL
	ABLSIQ
	ABLSMSKL
	ABLSMSKQ
	ABLSRL
	ABLSRQ
	ABOUNDL
	ABOUNDW
	ABSFL
	ABSFQ
	ABSFW
	ABSRL
	ABSRQ
	ABSRW
	ABSWAPL
	ABSWAPQ
	ABTCL
	ABTCQ
	ABTCW
	ABTL
	ABTQ
	ABTRL
	ABTRQ
	ABTRW
	ABTSL
	ABTSQ
	ABTSW
	ABTW
	ABYTE
	ABZHIL
	ABZHIQ
	ACBW
	ACDQ
	ACDQE
	ACLAC
	ACLC
	ACLD
	ACLDEMOTE
	ACLFLUSH
	ACLFLUSHOPT
	ACLI
	ACLTS
	ACLWB
	ACMC
	ACMOVLCC
	ACMOVLCS
	ACMOVLEQ
	ACMOVLGE
	ACMOVLGT
	ACMOVLHI
	ACMOVLLE
	ACMOVLLS
	ACMOVLLT
	ACMOVLMI
	ACMOVLNE
	ACMOVLOC
	ACMOVLOS
	ACMOVLPC
	ACMOVLPL
	ACMOVLPS
	ACMOVQCC
	ACMOVQCS
	ACMOVQEQ
	ACMOVQGE
	ACMOVQGT
	ACMOVQHI
	ACMOVQLE
	ACMOVQLS
	ACMOVQLT
	ACMOVQMI
	ACMOVQNE
	ACMOVQOC
	ACMOVQOS
	ACMOVQPC
	ACMOVQPL
	ACMOVQPS
	ACMOVWCC
	ACMOVWCS
	ACMOVWEQ
	ACMOVWGE
	ACMOVWGT
	ACMOVWHI
	ACMOVWLE
	ACMOVWLS
	ACMOVWLT
	ACMOVWMI
	ACMOVWNE
	ACMOVWOC
	ACMOVWOS
	ACMOVWPC
	ACMOVWPL
	ACMOVWPS
	ACMPB
	ACMPL
	ACMPPD
	ACMPPS
	ACMPQ
	ACMPSB
	ACMPSD
	ACMPSL
	ACMPSQ
	ACMPSS
	ACMPSW
	ACMPW
	ACMPXCHG16B
	ACMPXCHG8B
	ACMPXCHGB
	ACMPXCHGL
	ACMPXCHGQ
	ACMPXCHGW
	ACOMISD
	ACOMISS
	ACPUID
	ACQO
	ACRC32B
	ACRC32L
	ACRC32Q
	ACRC32W
	ACVTPD2PL
	ACVTPD2PS
	ACVTPL2PD
	ACVTPL2PS
	ACVTPS2PD
	ACVTPS2PL
	ACVTSD2SL
	ACVTSD2SQ
	ACVTSD2SS
	ACVTSL2SD
	ACVTSL2SS
	ACVTSQ2SD
	ACVTSQ2SS
	ACVTSS2SD
	ACVTSS2SL
	ACVTSS2SQ
	ACVTTPD2PL
	ACVTTPS2PL
	ACVTTSD2SL
	ACVTTSD2SQ
	ACVTTSS2SL
	ACVTTSS2SQ
	ACWD
	ACWDE
	ADAA
	ADAS
	ADECB
	ADECL
	ADECQ
	ADECW
	ADIVB
	ADIVL
	ADIVPD
	ADIVPS
	ADIVQ
	ADIVSD
	ADIVSS
	ADIVW
	ADPPD
	ADPPS
	AEMMS
	AENTER
	AEXTRACTPS
	AF2XM1
	AFABS
	AFADDD
	AFADDDP
	AFADDF
	AFADDL
	AFADDW
	AFBLD
	AFBSTP
	AFCHS
	AFCLEX
	AFCMOVB
	AFCMOVBE
	AFCMOVCC
	AFCMOVCS
	AFCMOVE
	AFCMOVEQ
	AFCMOVHI
	AFCMOVLS
	AFCMOVNB
	AFCMOVNBE
	AFCMOVNE
	AFCMOVNU
	AFCMOVU
	AFCMOVUN
	AFCOMD
	AFCOMDP
	AFCOMDPP
	AFCOMF
	AFCOMFP
	AFCOMI
	AFCOMIP
	AFCOML
	AFCOMLP
	AFCOMW
	AFCOMWP
	AFCOS
	AFDECSTP
	AFDIVD
	AFDIVDP
	AFDIVF
	AFDIVL
	AFDIVRD
	AFDIVRDP
	AFDIVRF
	AFDIVRL
	AFDIVRW
	AFDIVW
	AFFREE
	AFINCSTP
	AFINIT
	AFLD1
	AFLDCW
	AFLDENV
	AFLDL2E
	AFLDL2T
	AFLDLG2
	AFLDLN2
	AFLDPI
	AFLDZ
	AFMOVB
	AFMOVBP
	AFMOVD
	AFMOVDP
	AFMOVF
	AFMOVFP
	AFMOVL
	AFMOVLP
	AFMOVV
	AFMOVVP
	AFMOVW
	AFMOVWP
	AFMOVX
	AFMOVXP
	AFMULD
	AFMULDP
	AFMULF
	AFMULL
	AFMULW
	AFNOP
	AFPATAN
	AFPREM
	AFPREM1
	AFPTAN
	AFRNDINT
	AFRSTOR
	AFSAVE
	AFSCALE
	AFSIN
	AFSINCOS
	AFSQRT
	AFSTCW
	AFSTENV
	AFSTSW
	AFSUBD
	AFSUBDP
	AFSUBF
	AFSUBL
	AFSUBRD
	AFSUBRDP
	AFSUBRF
	AFSUBRL
	AFSUBRW
	AFSUBW
	AFTST
	AFUCOM
	AFUCOMI
	AFUCOMIP
	AFUCOMP
	AFUCOMPP
	AFXAM
	AFXCHD
	AFXRSTOR
	AFXRSTOR64
	AFXSAVE
	AFXSAVE64
	AFXTRACT
	AFYL2X
	AFYL2XP1
	AHADDPD
	AHADDPS
	AHLT
	AHSUBPD
	AHSUBPS
	AICEBP
	AIDIVB
	AIDIVL
	AIDIVQ
	AIDIVW
	AIMUL3L
	AIMUL3Q
	AIMUL3W
	AIMULB
	AIMULL
	AIMULQ
	AIMULW
	AINB
	AINCB
	AINCL
	AINCQ
	AINCW
	AINL
	AINSB
	AINSERTPS
	AINSL
	AINSW
	AINT
	AINTO
	AINVD
	AINVLPG
	AINVPCID
	AINW
	AIRETL
	AIRETQ
	AIRETW
	AJCC // >= unsigned
	AJCS // < unsigned
	AJCXZL
	AJCXZQ
	AJCXZW
	AJEQ // == (zero)
	AJGE // >= signed
	AJGT // > signed
	AJHI // > unsigned
	AJLE // <= signed
	AJLS // <= unsigned
	AJLT // < signed
	AJMI // sign bit set (negative)
	AJNE // != (nonzero)
	AJOC // overflow clear
	AJOS // overflow set
	AJPC // parity clear
	AJPL // sign bit clear (positive)
	AJPS // parity set
	AKADDB
	AKADDD
	AKADDQ
	AKADDW
	AKANDB
	AKANDD
	AKANDNB
	AKANDND
	AKANDNQ
	AKANDNW
	AKANDQ
	AKANDW
	AKMOVB
	AKMOVD
	AKMOVQ
	AKMOVW
	AKNOTB
	AKNOTD
	AKNOTQ
	AKNOTW
	AKORB
	AKORD
	AKORQ
	AKORTESTB
	AKORTESTD
	AKORTESTQ
	AKORTESTW
	AKORW
	AKSHIFTLB
	AKSHIFTLD
	AKSHIFTLQ
	AKSHIFTLW
	AKSHIFTRB
	AKSHIFTRD
	AKSHIFTRQ
	AKSHIFTRW
	AKTESTB
	AKTESTD
	AKTESTQ
	AKTESTW
	AKUNPCKBW
	AKUNPCKDQ
	AKUNPCKWD
	AKXNORB
	AKXNORD
	AKXNORQ
	AKXNORW
	AKXORB
	AKXORD
	AKXORQ
	AKXORW
	ALAHF
	ALARL
	ALARQ
	ALARW
	ALDDQU
	ALDMXCSR
	ALEAL
	ALEAQ
	ALEAVEL
	ALEAVEQ
	ALEAVEW
	ALEAW
	ALFENCE
	ALFSL
	ALFSQ
	ALFSW
	ALGDT
	ALGSL
	ALGSQ
	ALGSW
	ALIDT
	ALLDT
	ALMSW
	ALOCK
	ALODSB
	ALODSL
	ALODSQ
	ALODSW
	ALONG
	ALOOP
	ALOOPEQ
	ALOOPNE
	ALSLL
	ALSLQ
	ALSLW
	ALSSL
	ALSSQ
	ALSSW
	ALTR
	ALZCNTL
	ALZCNTQ
	ALZCNTW
	AMASKMOVOU
	AMASKMOVQ
	AMAXPD
	AMAXPS
	AMAXSD
	AMAXSS
	AMFENCE
	AMINPD
	AMINPS
	AMINSD
	AMINSS
	AMONITOR
	AMOVAPD
	AMOVAPS
	AMOVB
	AMOVBEL
	AMOVBEQ
	AMOVBEW
	AMOVBLSX
	AMOVBLZX
	AMOVBQSX
	AMOVBQZX
	AMOVBWSX
	AMOVBWZX
	AMOVDDUP
	AMOVHLPS
	AMOVHPD
	AMOVHPS
	AMOVL
	AMOVLHPS
	AMOVLPD
	AMOVLPS
	AMOVLQSX
	AMOVLQZX
	AMOVMSKPD
	AMOVMSKPS
	AMOVNTDQA
	AMOVNTIL
	AMOVNTIQ
	AMOVNTO
	AMOVNTPD
	AMOVNTPS
	AMOVNTQ
	AMOVO
	AMOVOU
	AMOVQ
	AMOVQL
	AMOVQOZX
	AMOVSB
	AMOVSD
	AMOVSHDUP
	AMOVSL
	AMOVSLDUP
	AMOVSQ
	AMOVSS
	AMOVSW
	AMOVSWW
	AMOVUPD
	AMOVUPS
	AMOVW
	AMOVWLSX
	AMOVWLZX
	AMOVWQSX
	AMOVWQZX
	AMOVZWW
	AMPSADBW
	AMULB
	AMULL
	AMULPD
	AMULPS
	AMULQ
	AMULSD
	AMULSS
	AMULW
	AMULXL
	AMULXQ
	AMWAIT
	ANEGB
	ANEGL
	ANEGQ
	ANEGW
	ANOPL
	ANOPW
	ANOTB
	ANOTL
	ANOTQ
	ANOTW
	AORB
	AORL
	AORPD
	AORPS
	AORQ
	AORW
	AOUTB
	AOUTL
	AOUTSB
	AOUTSL
	AOUTSW
	AOUTW
	APABSB
	APABSD
	APABSW
	APACKSSLW
	APACKSSWB
	APACKUSDW
	APACKUSWB
	APADDB
	APADDL
	APADDQ
	APADDSB
	APADDSW
	APADDUSB
	APADDUSW
	APADDW
	APALIGNR
	APAND
	APANDN
	APAUSE
	APAVGB
	APAVGW
	APBLENDVB
	APBLENDW
	APCLMULQDQ
	APCMPEQB
	APCMPEQL
	APCMPEQQ
	APCMPEQW
	APCMPESTRI
	APCMPESTRM
	APCMPGTB
	APCMPGTL
	APCMPGTQ
	APCMPGTW
	APCMPISTRI
	APCMPISTRM
	APDEPL
	APDEPQ
	APEXTL
	APEXTQ
	APEXTRB
	APEXTRD
	APEXTRQ
	APEXTRW
	APHADDD
	APHADDSW
	APHADDW
	APHMINPOSUW
	APHSUBD
	APHSUBSW
	APHSUBW
	APINSRB
	APINSRD
	APINSRQ
	APINSRW
	APMADDUBSW
	APMADDWL
	APMAXSB
	APMAXSD
	APMAXSW
	APMAXUB
	APMAXUD
	APMAXUW
	APMINSB
	APMINSD
	APMINSW
	APMINUB
	APMINUD
	APMINUW
	APMOVMSKB
	APMOVSXBD
	APMOVSXBQ
	APMOVSXBW
	APMOVSXDQ
	APMOVSXWD
	APMOVSXWQ
	APMOVZXBD
	APMOVZXBQ
	APMOVZXBW
	APMOVZXDQ
	APMOVZXWD
	APMOVZXWQ
	APMULDQ
	APMULHRSW
	APMULHUW
	APMULHW
	APMULLD
	APMULLW
	APMULULQ
	APOPAL
	APOPAW
	APOPCNTL
	APOPCNTQ
	APOPCNTW
	APOPFL
	APOPFQ
	APOPFW
	APOPL
	APOPQ
	APOPW
	APOR
	APREFETCHNTA
	APREFETCHT0
	APREFETCHT1
	APREFETCHT2
	APSADBW
	APSHUFB
	APSHUFD
	APSHUFHW
	APSHUFL
	APSHUFLW
	APSHUFW
	APSIGNB
	APSIGND
	APSIGNW
	APSLLL
	APSLLO
	APSLLQ
	APSLLW
	APSRAL
	APSRAW
	APSRLL
	APSRLO
	APSRLQ
	APSRLW
	APSUBB
	APSUBL
	APSUBQ
	APSUBSB
	APSUBSW
	APSUBUSB
	APSUBUSW
	APSUBW
	APTEST
	APUNPCKHBW
	APUNPCKHLQ
	APUNPCKHQDQ
	APUNPCKHWL
	APUNPCKLBW
	APUNPCKLLQ
	APUNPCKLQDQ
	APUNPCKLWL
	APUSHAL
	APUSHAW
	APUSHFL
	APUSHFQ
	APUSHFW
	APUSHL
	APUSHQ
	APUSHW
	APXOR
	AQUAD
	ARCLB
	ARCLL
	ARCLQ
	ARCLW
	ARCPPS
	ARCPSS
	ARCRB
	ARCRL
	ARCRQ
	ARCRW
	ARDFSBASEL
	ARDFSBASEQ
	ARDGSBASEL
	ARDGSBASEQ
	ARDMSR
	ARDPKRU
	ARDPMC
	ARDRANDL
	ARDRANDQ
	ARDRANDW
	ARDSEEDL
	ARDSEEDQ
	ARDSEEDW
	ARDTSC
	ARDTSCP
	AREP
	AREPN
	ARETFL
	ARETFQ
	ARETFW
	AROLB
	AROLL
	AROLQ
	AROLW
	ARORB
	ARORL
	ARORQ
	ARORW
	ARORXL
	ARORXQ
	AROUNDPD
	AROUNDPS
	AROUNDSD
	AROUNDSS
	ARSM
	ARSQRTPS
	ARSQRTSS
	ASAHF
	ASALB
	ASALL
	ASALQ
	ASALW
	ASARB
	ASARL
	ASARQ
	ASARW
	ASARXL
	ASARXQ
	ASBBB
	ASBBL
	ASBBQ
	ASBBW
	ASCASB
	ASCASL
	ASCASQ
	ASCASW
	ASETCC
	ASETCS
	ASETEQ
	ASETGE
	ASETGT
	ASETHI
	ASETLE
	ASETLS
	ASETLT
	ASETMI
	ASETNE
	ASETOC
	ASETOS
	ASETPC
	ASETPL
	ASETPS
	ASFENCE
	ASGDT
	ASHA1MSG1
	ASHA1MSG2
	ASHA1NEXTE
	ASHA1RNDS4
	ASHA256MSG1
	ASHA256MSG2
	ASHA256RNDS2
	ASHLB
	ASHLL
	ASHLQ
	ASHLW
	ASHLXL
	ASHLXQ
	ASHRB
	ASHRL
	ASHRQ
	ASHRW
	ASHRXL
	ASHRXQ
	ASHUFPD
	ASHUFPS
	ASIDT
	ASLDTL
	ASLDTQ
	ASLDTW
	ASMSWL
	ASMSWQ
	ASMSWW
	ASQRTPD
	ASQRTPS
	ASQRTSD
	ASQRTSS
	ASTAC
	ASTC
	ASTD
	ASTI
	ASTMXCSR
	ASTOSB
	ASTOSL
	ASTOSQ
	ASTOSW
	ASTRL
	ASTRQ
	ASTRW
	ASUBB
	ASUBL
	ASUBPD
	ASUBPS
	ASUBQ
	ASUBSD
	ASUBSS
	ASUBW
	ASWAPGS
	ASYSCALL
	ASYSENTER
	ASYSENTER64
	ASYSEXIT
	ASYSEXIT64
	ASYSRET
	ATESTB
	ATESTL
	ATESTQ
	ATESTW
	ATPAUSE
	ATZCNTL
	ATZCNTQ
	ATZCNTW
	AUCOMISD
	AUCOMISS
	AUD1
	AUD2
	AUMWAIT
	AUNPCKHPD
	AUNPCKHPS
	AUNPCKLPD
	AUNPCKLPS
	AUMONITOR
	AV4FMADDPS
	AV4FMADDSS
	AV4FNMADDPS
	AV4FNMADDSS
	AVADDPD
	AVADDPS
	AVADDSD
	AVADDSS
	AVADDSUBPD
	AVADDSUBPS
	AVAESDEC
	AVAESDECLAST
	AVAESENC
	AVAESENCLAST
	AVAESIMC
	AVAESKEYGENASSIST
	AVALIGND
	AVALIGNQ
	AVANDNPD
	AVANDNPS
	AVANDPD
	AVANDPS
	AVBLENDMPD
	AVBLENDMPS
	AVBLENDPD
	AVBLENDPS
	AVBLENDVPD
	AVBLENDVPS
	AVBROADCASTF128
	AVBROADCASTF32X2
	AVBROADCASTF32X4
	AVBROADCASTF32X8
	AVBROADCASTF64X2
	AVBROADCASTF64X4
	AVBROADCASTI128
	AVBROADCASTI32X2
	AVBROADCASTI32X4
	AVBROADCASTI32X8
	AVBROADCASTI64X2
	AVBROADCASTI64X4
	AVBROADCASTSD
	AVBROADCASTSS
	AVCMPPD
	AVCMPPS
	AVCMPSD
	AVCMPSS
	AVCOMISD
	AVCOMISS
	AVCOMPRESSPD
	AVCOMPRESSPS
	AVCVTDQ2PD
	AVCVTDQ2PS
	AVCVTPD2DQ
	AVCVTPD2DQX
	AVCVTPD2DQY
	AVCVTPD2PS
	AVCVTPD2PSX
	AVCVTPD2PSY
	AVCVTPD2QQ
	AVCVTPD2UDQ
	AVCVTPD2UDQX
	AVCVTPD2UDQY
	AVCVTPD2UQQ
	AVCVTPH2PS
	AVCVTPS2DQ
	AVCVTPS2PD
	AVCVTPS2PH
	AVCVTPS2QQ
	AVCVTPS2UDQ
	AVCVTPS2UQQ
	AVCVTQQ2PD
	AVCVTQQ2PS
	AVCVTQQ2PSX
	AVCVTQQ2PSY
	AVCVTSD2SI
	AVCVTSD2SIQ
	AVCVTSD2SS
	AVCVTSD2USI
	AVCVTSD2USIL
	AVCVTSD2USIQ
	AVCVTSI2SDL
	AVCVTSI2SDQ
	AVCVTSI2SSL
	AVCVTSI2SSQ
	AVCVTSS2SD
	AVCVTSS2SI
	AVCVTSS2SIQ
	AVCVTSS2USI
	AVCVTSS2USIL
	AVCVTSS2USIQ
	AVCVTTPD2DQ
	AVCVTTPD2DQX
	AVCVTTPD2DQY
	AVCVTTPD2QQ
	AVCVTTPD2UDQ
	AVCVTTPD2UDQX
	AVCVTTPD2UDQY
	AVCVTTPD2UQQ
	AVCVTTPS2DQ
	AVCVTTPS2QQ
	AVCVTTPS2UDQ
	AVCVTTPS2UQQ
	AVCVTTSD2SI
	AVCVTTSD2SIQ
	AVCVTTSD2USI
	AVCVTTSD2USIL
	AVCVTTSD2USIQ
	AVCVTTSS2SI
	AVCVTTSS2SIQ
	AVCVTTSS2USI
	AVCVTTSS2USIL
	AVCVTTSS2USIQ
	AVCVTUDQ2PD
	AVCVTUDQ2PS
	AVCVTUQQ2PD
	AVCVTUQQ2PS
	AVCVTUQQ2PSX
	AVCVTUQQ2PSY
	AVCVTUSI2SD
	AVCVTUSI2SDL
	AVCVTUSI2SDQ
	AVCVTUSI2SS
	AVCVTUSI2SSL
	AVCVTUSI2SSQ
	AVDBPSADBW
	AVDIVPD
	AVDIVPS
	AVDIVSD
	AVDIVSS
	AVDPPD
	AVDPPS
	AVERR
	AVERW
	AVEXP2PD
	AVEXP2PS
	AVEXPANDPD
	AVEXPANDPS
	AVEXTRACTF128
	AVEXTRACTF32X4
	AVEXTRACTF32X8
	AVEXTRACTF64X2
	AVEXTRACTF64X4
	AVEXTRACTI128
	AVEXTRACTI32X4
	AVEXTRACTI32X8
	AVEXTRACTI64X2
	AVEXTRACTI64X4
	AVEXTRACTPS
	AVFIXUPIMMPD
	AVFIXUPIMMPS
	AVFIXUPIMMSD
	AVFIXUPIMMSS
	AVFMADD132PD
	AVFMADD132PS
	AVFMADD132SD
	AVFMADD132SS
	AVFMADD213PD
	AVFMADD213PS
	AVFMADD213SD
	AVFMADD213SS
	AVFMADD231PD
	AVFMADD231PS
	AVFMADD231SD
	AVFMADD231SS
	AVFMADDSUB132PD
	AVFMADDSUB132PS
	AVFMADDSUB213PD
	AVFMADDSUB213PS
	AVFMADDSUB231PD
	AVFMADDSUB231PS
	AVFMSUB132PD
	AVFMSUB132PS
	AVFMSUB132SD
	AVFMSUB132SS
	AVFMSUB213PD
	AVFMSUB213PS
	AVFMSUB213SD
	AVFMSUB213SS
	AVFMSUB231PD
	AVFMSUB231PS
	AVFMSUB231SD
	AVFMSUB231SS
	AVFMSUBADD132PD
	AVFMSUBADD132PS
	AVFMSUBADD213PD
	AVFMSUBADD213PS
	AVFMSUBADD231PD
	AVFMSUBADD231PS
	AVFNMADD132PD
	AVFNMADD132PS
	AVFNMADD132SD
	AVFNMADD132SS
	AVFNMADD213PD
	AVFNMADD213PS
	AVFNMADD213SD
	AVFNMADD213SS
	AVFNMADD231PD
	AVFNMADD231PS
	AVFNMADD231SD
	AVFNMADD231SS
	AVFNMSUB132PD
	AVFNMSUB132PS
	AVFNMSUB132SD
	AVFNMSUB132SS
	AVFNMSUB213PD
	AVFNMSUB213PS
	AVFNMSUB213SD
	AVFNMSUB213SS
	AVFNMSUB231PD
	AVFNMSUB231PS
	AVFNMSUB231SD
	AVFNMSUB231SS
	AVFPCLASSPD
	AVFPCLASSPDX
	AVFPCLASSPDY
	AVFPCLASSPDZ
	AVFPCLASSPS
	AVFPCLASSPSX
	AVFPCLASSPSY
	AVFPCLASSPSZ
	AVFPCLASSSD
	AVFPCLASSSS
	AVGATHERDPD
	AVGATHERDPS
	AVGATHERPF0DPD
	AVGATHERPF0DPS
	AVGATHERPF0QPD
	AVGATHERPF0QPS
	AVGATHERPF1DPD
	AVGATHERPF1DPS
	AVGATHERPF1QPD
	AVGATHERPF1QPS
	AVGATHERQPD
	AVGATHERQPS
	AVGETEXPPD
	AVGETEXPPS
	AVGETEXPSD
	AVGETEXPSS
	AVGETMANTPD
	AVGETMANTPS
	AVGETMANTSD
	AVGETMANTSS
	AVGF2P8AFFINEINVQB
	AVGF2P8AFFINEQB
	AVGF2P8MULB
	AVHADDPD
	AVHADDPS
	AVHSUBPD
	AVHSUBPS
	AVINSERTF128
	AVINSERTF32X4
	AVINSERTF32X8
	AVINSERTF64X2
	AVINSERTF64X4
	AVINSERTI128
	AVINSERTI32X4
	AVINSERTI32X8
	AVINSERTI64X2
	AVINSERTI64X4
	AVINSERTPS
	AVLDDQU
	AVLDMXCSR
	AVMASKMOVDQU
	AVMASKMOVPD
	AVMASKMOVPS
	AVMAXPD
	AVMAXPS
	AVMAXSD
	AVMAXSS
	AVMINPD
	AVMINPS
	AVMINSD
	AVMINSS
	AVMOVAPD
	AVMOVAPS
	AVMOVD
	AVMOVDDUP
	AVMOVDQA
	AVMOVDQA32
	AVMOVDQA64
	AVMOVDQU
	AVMOVDQU16
	AVMOVDQU32
	AVMOVDQU64
	AVMOVDQU8
	AVMOVHLPS
	AVMOVHPD
	AVMOVHPS
	AVMOVLHPS
	AVMOVLPD
	AVMOVLPS
	AVMOVMSKPD
	AVMOVMSKPS
	AVMOVNTDQ
	AVMOVNTDQA
	AVMOVNTPD
	AVMOVNTPS
	AVMOVQ
	AVMOVSD
	AVMOVSHDUP
	AVMOVSLDUP
	AVMOVSS
	AVMOVUPD
	AVMOVUPS
	AVMPSADBW
	AVMULPD
	AVMULPS
	AVMULSD
	AVMULSS
	AVORPD
	AVORPS
	AVP4DPWSSD
	AVP4DPWSSDS
	AVPABSB
	AVPABSD
	AVPABSQ
	AVPABSW
	AVPACKSSDW
	AVPACKSSWB
	AVPACKUSDW
	AVPACKUSWB
	AVPADDB
	AVPADDD
	AVPADDQ
	AVPADDSB
	AVPADDSW
	AVPADDUSB
	AVPADDUSW
	AVPADDW
	AVPALIGNR
	AVPAND
	AVPANDD
	AVPANDN
	AVPANDND
	AVPANDNQ
	AVPANDQ
	AVPAVGB
	AVPAVGW
	AVPBLENDD
	AVPBLENDMB
	AVPBLENDMD
	AVPBLENDMQ
	AVPBLENDMW
	AVPBLENDVB
	AVPBLENDW
	AVPBROADCASTB
	AVPBROADCASTD
	AVPBROADCASTMB2Q
	AVPBROADCASTMW2D
	AVPBROADCASTQ
	AVPBROADCASTW
	AVPCLMULQDQ
	AVPCMPB
	AVPCMPD
	AVPCMPEQB
	AVPCMPEQD
	AVPCMPEQQ
	AVPCMPEQW
	AVPCMPESTRI
	AVPCMPESTRM
	AVPCMPGTB
	AVPCMPGTD
	AVPCMPGTQ
	AVPCMPGTW
	AVPCMPISTRI
	AVPCMPISTRM
	AVPCMPQ
	AVPCMPUB
	AVPCMPUD
	AVPCMPUQ
	AVPCMPUW
	AVPCMPW
	AVPCOMPRESSB
	AVPCOMPRESSD
	AVPCOMPRESSQ
	AVPCOMPRESSW
	AVPCONFLICTD
	AVPCONFLICTQ
	AVPDPBUSD
	AVPDPBUSDS
	AVPDPWSSD
	AVPDPWSSDS
	AVPERM2F128
	AVPERM2I128
	AVPERMB
	AVPERMD
	AVPERMI2B
	AVPERMI2D
	AVPERMI2PD
	AVPERMI2PS
	AVPERMI2Q
	AVPERMI2W
	AVPERMILPD
	AVPERMILPS
	AVPERMPD
	AVPERMPS
	AVPERMQ
	AVPERMT2B
	AVPERMT2D
	AVPERMT2PD
	AVPERMT2PS
	AVPERMT2Q
	AVPERMT2W
	AVPERMW
	AVPEXPANDB
	AVPEXPANDD
	AVPEXPANDQ
	AVPEXPANDW
	AVPEXTRB
	AVPEXTRD
	AVPEXTRQ
	AVPEXTRW
	AVPGATHERDD
	AVPGATHERDQ
	AVPGATHERQD
	AVPGATHERQQ
	AVPHADDD
	AVPHADDSW
	AVPHADDW
	AVPHMINPOSUW
	AVPHSUBD
	AVPHSUBSW
	AVPHSUBW
	AVPINSRB
	AVPINSRD
	AVPINSRQ
	AVPINSRW
	AVPLZCNTD
	AVPLZCNTQ
	AVPMADD52HUQ
	AVPMADD52LUQ
	AVPMADDUBSW
	AVPMADDWD
	AVPMASKMOVD
	AVPMASKMOVQ
	AVPMAXSB
	AVPMAXSD
	AVPMAXSQ
	AVPMAXSW
	AVPMAXUB
	AVPMAXUD
	AVPMAXUQ
	AVPMAXUW
	AVPMINSB
	AVPMINSD
	AVPMINSQ
	AVPMINSW
	AVPMINUB
	AVPMINUD
	AVPMINUQ
	AVPMINUW
	AVPMOVB2M
	AVPMOVD2M
	AVPMOVDB
	AVPMOVDW
	AVPMOVM2B
	AVPMOVM2D
	AVPMOVM2Q
	AVPMOVM2W
	AVPMOVMSKB
	AVPMOVQ2M
	AVPMOVQB
	AVPMOVQD
	AVPMOVQW
	AVPMOVSDB
	AVPMOVSDW
	AVPMOVSQB
	AVPMOVSQD
	AVPMOVSQW
	AVPMOVSWB
	AVPMOVSXBD
	AVPMOVSXBQ
	AVPMOVSXBW
	AVPMOVSXDQ
	AVPMOVSXWD
	AVPMOVSXWQ
	AVPMOVUSDB
	AVPMOVUSDW
	AVPMOVUSQB
	AVPMOVUSQD
	AVPMOVUSQW
	AVPMOVUSWB
	AVPMOVW2M
	AVPMOVWB
	AVPMOVZXBD
	AVPMOVZXBQ
	AVPMOVZXBW
	AVPMOVZXDQ
	AVPMOVZXWD
	AVPMOVZXWQ
	AVPMULDQ
	AVPMULHRSW
	AVPMULHUW
	AVPMULHW
	AVPMULLD
	AVPMULLQ
	AVPMULLW
	AVPMULTISHIFTQB
	AVPMULUDQ
	AVPOPCNTB
	AVPOPCNTD
	AVPOPCNTQ
	AVPOPCNTW
	AVPOR
	AVPORD
	AVPORQ
	AVPROLD
	AVPROLQ
	AVPROLVD
	AVPROLVQ
	AVPRORD
	AVPRORQ
	AVPRORVD
	AVPRORVQ
	AVPSADBW
	AVPSCATTERDD
	AVPSCATTERDQ
	AVPSCATTERQD
	AVPSCATTERQQ
	AVPSHLDD
	AVPSHLDQ
	AVPSHLDVD
	AVPSHLDVQ
	AVPSHLDVW
	AVPSHLDW
	AVPSHRDD
	AVPSHRDQ
	AVPSHRDVD
	AVPSHRDVQ
	AVPSHRDVW
	AVPSHRDW
	AVPSHUFB
	AVPSHUFBITQMB
	AVPSHUFD
	AVPSHUFHW
	AVPSHUFLW
	AVPSIGNB
	AVPSIGND
	AVPSIGNW
	AVPSLLD
	AVPSLLDQ
	AVPSLLQ
	AVPSLLVD
	AVPSLLVQ
	AVPSLLVW
	AVPSLLW
	AVPSRAD
	AVPSRAQ
	AVPSRAVD
	AVPSRAVQ
	AVPSRAVW
	AVPSRAW
	AVPSRLD
	AVPSRLDQ
	AVPSRLQ
	AVPSRLVD
	AVPSRLVQ
	AVPSRLVW
	AVPSRLW
	AVPSUBB
	AVPSUBD
	AVPSUBQ
	AVPSUBSB
	AVPSUBSW
	AVPSUBUSB
	AVPSUBUSW
	AVPSUBW
	AVPTERNLOGD
	AVPTERNLOGQ
	AVPTEST
	AVPTESTMB
	AVPTESTMD
	AVPTESTMQ
	AVPTESTMW
	AVPTESTNMB
	AVPTESTNMD
	AVPTESTNMQ
	AVPTESTNMW
	AVPUNPCKHBW
	AVPUNPCKHDQ
	AVPUNPCKHQDQ
	AVPUNPCKHWD
	AVPUNPCKLBW
	AVPUNPCKLDQ
	AVPUNPCKLQDQ
	AVPUNPCKLWD
	AVPXOR
	AVPXORD
	AVPXORQ
	AVRANGEPD
	AVRANGEPS
	AVRANGESD
	AVRANGESS
	AVRCP14PD
	AVRCP14PS
	AVRCP14SD
	AVRCP14SS
	AVRCP28PD
	AVRCP28PS
	AVRCP28SD
	AVRCP28SS
	AVRCPPS
	AVRCPSS
	AVREDUCEPD
	AVREDUCEPS
	AVREDUCESD
	AVREDUCESS
	AVRNDSCALEPD
	AVRNDSCALEPS
	AVRNDSCALESD
	AVRNDSCALESS
	AVROUNDPD
	AVROUNDPS
	AVROUNDSD
	AVROUNDSS
	AVRSQRT14PD
	AVRSQRT14PS
	AVRSQRT14SD
	AVRSQRT14SS
	AVRSQRT28PD
	AVRSQRT28PS
	AVRSQRT28SD
	AVRSQRT28SS
	AVRSQRTPS
	AVRSQRTSS
	AVSCALEFPD
	AVSCALEFPS
	AVSCALEFSD
	AVSCALEFSS
	AVSCATTERDPD
	AVSCATTERDPS
	AVSCATTERPF0DPD
	AVSCATTERPF0DPS
	AVSCATTERPF0QPD
	AVSCATTERPF0QPS
	AVSCATTERPF1DPD
	AVSCATTERPF1DPS
	AVSCATTERPF1QPD
	AVSCATTERPF1QPS
	AVSCATTERQPD
	AVSCATTERQPS
	AVSHUFF32X4
	AVSHUFF64X2
	AVSHUFI32X4
	AVSHUFI64X2
	AVSHUFPD
	AVSHUFPS
	AVSQRTPD
	AVSQRTPS
	AVSQRTSD
	AVSQRTSS
	AVSTMXCSR
	AVSUBPD
	AVSUBPS
	AVSUBSD
	AVSUBSS
	AVTESTPD
	AVTESTPS
	AVUCOMISD
	AVUCOMISS
	AVUNPCKHPD
	AVUNPCKHPS
	AVUNPCKLPD
	AVUNPCKLPS
	AVXORPD
	AVXORPS
	AVZEROALL
	AVZEROUPPER
	AWAIT
	AWBINVD
	AWORD
	AWRFSBASEL
	AWRFSBASEQ
	AWRGSBASEL
	AWRGSBASEQ
	AWRMSR
	AWRPKRU
	AXABORT
	AXACQUIRE
	AXADDB
	AXADDL
	AXADDQ
	AXADDW
	AXBEGIN
	AXCHGB
	AXCHGL
	AXCHGQ
	AXCHGW
	AXEND
	AXGETBV
	AXLAT
	AXORB
	AXORL
	AXORPD
	AXORPS
	AXORQ
	AXORW
	AXRELEASE
	AXRSTOR
	AXRSTOR64
	AXRSTORS
	AXRSTORS64
	AXSAVE
	AXSAVE64
	AXSAVEC
	AXSAVEC64
	AXSAVEOPT
	AXSAVEOPT64
	AXSAVES
	AXSAVES64
	AXSETBV
	AXTEST
	ALAST
)
