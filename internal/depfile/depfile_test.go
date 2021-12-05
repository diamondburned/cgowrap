package depfile

import (
	"reflect"
	"strings"
	"testing"

	_ "embed"
)

//go:embed depfile_test_data
var testdata string

func TestParseFile(t *testing.T) {
	f, err := ParseFile(strings.NewReader(testdata))
	if err != nil {
		t.Fatal("cannot parse file:", err)
	}

	expect := map[string][]string{
		"_obj/_7_cgo_.o": {
			"/tmp/cgo-gcc-input-2620350145.c",
			"/nix/store/cn4z6y3pzcr7pry9078rsmd81b8zg3y5-clang-wrapper-7.1.0/resource-root/include/stddef.h",
			"/nix/store/cn4z6y3pzcr7pry9078rsmd81b8zg3y5-clang-wrapper-7.1.0/resource-root/include/__stddef_max_align_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/stdlib.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/libc-header-start.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/features.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/stdc-predef.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/sys/cdefs.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/wordsize.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/long-double.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/gnu/stubs.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/gnu/stubs-64.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/waitflags.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/waitstatus.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/floatn.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/floatn-common.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/sys/types.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/timesize.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/typesizes.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/time64.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/clock_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/clockid_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/time_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/timer_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/stdint-intn.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/endian.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/endian.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/endianness.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/byteswap.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/uintn-identity.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/sys/select.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/select.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/sigset_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/__sigset_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/struct_timeval.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/struct_timespec.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/pthreadtypes.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/thread-shared-types.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/pthreadtypes-arch.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/struct_mutex.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/struct_rwlock.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/alloca.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/stdlib-float.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/galloca.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gtypes.h",
			"/nix/store/dcb3cyba5wl6qimv6vwdpbi0kg0g1nlb-glib-2.68.2/lib/glib-2.0/include/glibconfig.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gmacros.h",
			"/nix/store/cn4z6y3pzcr7pry9078rsmd81b8zg3y5-clang-wrapper-7.1.0/resource-root/include/limits.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/limits.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/posix1_lim.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/local_lim.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/linux/limits.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/posix2_lim.h",
			"/nix/store/cn4z6y3pzcr7pry9078rsmd81b8zg3y5-clang-wrapper-7.1.0/resource-root/include/float.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gversionmacros.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/time.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/time.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/struct_tm.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/struct_itimerspec.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/locale_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/__locale_t.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/garray.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gasyncqueue.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gthread.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gatomic.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gerror.h",
			"/nix/store/cn4z6y3pzcr7pry9078rsmd81b8zg3y5-clang-wrapper-7.1.0/resource-root/include/stdarg.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gquark.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gutils.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gbacktrace.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/signal.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/signum-generic.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/signum-arch.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/sig_atomic_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/siginfo_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/__sigval_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/siginfo-arch.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/siginfo-consts.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/sigval_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/sigevent_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/sigevent-consts.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/sigaction.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/sigcontext.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/stack_t.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/sys/ucontext.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/sigstack.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/ss_flags.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/struct_sigstack.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/sigthread.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/signal_ext.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gbase64.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gbitlock.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gbookmarkfile.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gdatetime.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gtimezone.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gbytes.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gcharset.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gchecksum.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gconvert.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gdataset.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gdate.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gdir.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/dirent.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/dirent.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/dirent_ext.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/genviron.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gfileutils.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/ggettext.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/ghash.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/glist.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gmem.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gnode.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/ghmac.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/ghook.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/ghostutils.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/giochannel.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gmain.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gpoll.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gslist.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gstring.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gunicode.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gkeyfile.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gmappedfile.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gmarkup.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gmessages.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gvariant.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gvarianttype.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/goption.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gpattern.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gprimes.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gqsort.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gqueue.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/grand.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/grcbox.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/grefcount.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/grefstring.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gregex.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gscanner.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gsequence.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gshell.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gslice.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/string.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/strings.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gspawn.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gstrfuncs.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gstringchunk.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gstrvbuilder.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gtestutils.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/errno.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/errno.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/linux/errno.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/asm/errno.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/asm-generic/errno.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/asm-generic/errno-base.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gthreadpool.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gtimer.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gtrashstack.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gtree.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/guri.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/guuid.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/gversion.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/deprecated/gallocator.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/deprecated/gcache.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/deprecated/gcompletion.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/deprecated/gmain.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/deprecated/grel.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/deprecated/gthread.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/pthread.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/sched.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/sched.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/types/struct_sched_param.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/cpu-set.h",
			"/nix/store/am5qwbpriqhp1i9qhp2idid7ympxqb9a-glibc-2.32-46-dev/include/bits/setjmp.h",
			"/nix/store/d9zs9xg86lhqjqni0v8h2ibdrjb57fn4-glib-2.68.2-dev/include/glib-2.0/glib/glib-autocleanups.h",
		},
	}

	if !reflect.DeepEqual(expect, f.Sources) {
		t.Errorf("expect: %#q", expect)
		t.Errorf("got:    %#q", f.Sources)
	}
}