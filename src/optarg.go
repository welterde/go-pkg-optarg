/*

 author: jim teeuwen <jimteeuwen@gmail.com>
 version: 0.1

 OPTARG - A simple commandline options parser.

 - Allows options with single and multi character names according to the
   traditional unix way of doing things. eg:  -o versus --option
 - Exposes a channel based iterator which returns parsed options from the
   optarg.Parse() function. Note that it only yields Options which are actually
   present in the commandline arguments (os.Args). Call as a for loop:

   for opt := range optarg.Parse() {
   	// ... parse option
   }

 - Standard switch tokens are - and --. Can be modified by changing the vars
   optarg.ShortSwitch and optarg.LongSwitch.
 - Any arguments not associated with an option will be available in the
   optarg.Remainder slice after optarg.Parse() has been run.
 - Boolean flags require no value. The precense or absence of the flag is the
   value by itself. eg: flag '-n' is false if it's not found in os.Args, true if
   it is.
 - Exposes a Usage() function which prints options with their description and
   default values to the standard output. As opposed to the flag package, this
   outputs *neatly formatted* text. It prints output like the following. Note
   that this is the Usage() output of the options listed in optarg_test.go. It
   uses my own sexy multilineWrap() routine (see string.go):

 Usage: ./6.out [options]:

 --source, -s: Path to the source folder. Here is some added description
               information which is completely useless, but it makes sure we can
               pimp our sexy Usage() output when dealing with lenghty, multi
               -line description texts.
    --bin, -b: Path to the binary folder.
   --arch, -a: Target architecture. (defaults to: amd64)
 --noproc, -n: Skip pre/post processing. (defaults to: false)
  --purge, -p: Clean compiled packages after linking is complete. (defaults to:
               false)
 
*/
package optarg

import "fmt"
import "os"
import "strings"
import "strconv"

type Option struct {
	Name		string;
	ShortName	string;
	Description	string;

	defaultval	interface{};
	value		string;
}

var (
	options		= make([]*Option, 0);
	Remainder	= make([]string, 0);
	ShortSwitch	= "-";
	LongSwitch	= "--";
	appname		= os.Args[0];
)

// Prints usage information in a neatly formatted overview.
func Usage() {
	offset := 0;

	// Find the largest length of the option name list. Needed to align
	// the description blocks consistently.
	for _, v := range options {
		str := fmt.Sprintf("%s%s, %s%s: ", LongSwitch, v.Name, ShortSwitch, v.ShortName);
		if len(str) > offset {
			offset = len(str);
		}
	}

	offset++; // add margin.

	fmt.Printf("Usage: %s [options]:\n\n", appname);

	for _, v := range options {
		// Print namelist. right-align it based on the maximum width
		// found in previous loop.
		str := fmt.Sprintf("%s%s, %s%s: ", LongSwitch, v.Name, ShortSwitch, v.ShortName);
		format := fmt.Sprintf("%%%ds", offset);
		fmt.Printf(format, str);

		desc := v.Description;
		if fmt.Sprintf("%v", v.defaultval) != "" {
			desc = fmt.Sprintf("%s (defaults to: %v)", desc, v.defaultval);
		}

		// Format and print left-aligned, word-wrapped description with
		// a @margin left margin size using my super duper patented
		// multi-line string wrap routine (see string.go). Assume
		// maximum of 80 characters screen width. Which makes block
		// width equal to 80 - @offset. I would prefer to use
		// ALIGN_JUSTIFY for added sexy, but it looks a little funky for
		// short descriptions. So we'll stick with the establish left-
		// aligned text.
		lines := multilineWrap(desc, 80, offset, 0, ALIGN_LEFT);
		
		// First line needs to be appended to where we left off.
		fmt.Printf("%s\n", strings.TrimSpace(lines[0]));
		
		// Print the rest as-is (properly indented).
		for i := 1; i < len(lines); i++ {
			fmt.Printf("%s\n", lines[i]);
		}
	}

	offset++; // add margin.
}

// Parse os.Args using the previously added Options.
func Parse() <-chan *Option {
	c := make(chan *Option);
	Remainder = make([]string, 0);
	go processArgs(c);
	return c;
}

func processArgs(c chan *Option) {
	var opt *Option;
	for i, v := range os.Args {
		if i == 0 {
			continue
		} // skip app name
	
		v := strings.TrimSpace(v);
		if len(v) == 0 {
			continue
		}

		if len(v) >= 3 && v[0:2] == LongSwitch {
			v := strings.TrimSpace(v[2:len(v)]);
			if len(v) == 0 {
				listAppend(&Remainder, LongSwitch)
			} else {
				opt = findOption(v);
				if opt == nil {
					fmt.Fprintf(os.Stderr, "Unknown option '--%s' specified.\n", v);
					Usage();
					os.Exit(1);
				}

				_, ok := opt.defaultval.(bool);
				if ok {
					opt.value = "true";
					c <- opt;
					opt = nil;
				}
			}

		} else if len(v) >= 2 && v[0:1] == ShortSwitch {
			v := strings.TrimSpace(v[1:len(v)]);
			if len(v) == 0 {
				listAppend(&Remainder, ShortSwitch)
			} else {
				for i, _ := range v {
					tok := v[i : i+1];
					opt = findOption(tok);
					if opt == nil {
						fmt.Fprintf(os.Stderr, "Unknown option '-%s' specified.\n", tok);
						Usage();
						os.Exit(1);
					}

					_, ok := opt.defaultval.(bool);
					if ok {
						opt.value = "true";
						c <- opt;
						opt = nil;
					}
				}
			}

		} else {
			if opt == nil {
				listAppend(&Remainder, v)
			} else {
				opt.value = v;
				c <- opt;
				opt = nil;
			}
		}
	}
	close(c);
}

// Add a new command line option to check for.
func Add(shortname, name, description string, defaultvalue interface{}) {
	opt := &Option{
		ShortName: shortname,
		Name: name,
		Description: description,
		defaultval: defaultvalue,
	};

	slice := make([]*Option, len(options)+1);
	for i, v := range options {
		slice[i] = v
	}
	slice[len(slice)-1] = opt;
	options = slice;
}

func findOption(name string) *Option {
	for _, opt := range options {
		if opt.Name == name || opt.ShortName == name {
			return opt
		}
	}
	return nil;
}

func (this *Option) String() string	{ return this.value }

func (this *Option) Bool() bool {
	yes := []string{"1", "y", "yes", "true", "on"};
	this.value = strings.ToLower(this.value);
	for _, v := range yes {
		if v == this.value {
			return true
		}
	}
	return false;
}

func (this *Option) Int() int {
	v, err := strconv.Atoi(this.value);
	if err != nil {
		return this.defaultval.(int)
	}
	return v;
}

func (this *Option) Int64() int64 {
	v, err := strconv.Atoi64(this.value);
	if err != nil {
		return this.defaultval.(int64)
	}
	return v;
}

func (this *Option) Float() float {
	v, err := strconv.Atof(this.value);
	if err != nil {
		return this.defaultval.(float)
	}
	return v;
}

func (this *Option) Float32() float32 {
	v, err := strconv.Atof32(this.value);
	if err != nil {
		return this.defaultval.(float32)
	}
	return v;
}

func (this *Option) Float64() float64 {
	v, err := strconv.Atof64(this.value);
	if err != nil {
		return this.defaultval.(float64)
	}
	return v;
}

func listAppend(list *[]string, item string) {
	slice := make([]string, len(*list)+1);
	for i, v := range *list {
		slice[i] = v
	}
	slice[len(slice)-1] = item;
	*list = slice;
}
