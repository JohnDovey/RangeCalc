# GO Language version

This Directory holds the *Go* language version of the Range Calculator

- The file `RangeCalc.exe` is the Windows compiled version. 
- `RangeCalc` was compiled on Ubuntu 20.04.
- Source is, of course, in `RangeCalc.go`.

## Usage

Run with no flags for the original interactive prompts. Pass `-b1`, `-b2`,
and `-d` together to skip straight to a calculation — handy for scripting
or piping:

```
RangeCalc -b1 25 -b2 350 -d 250
RangeFromRef1=357.04
RangeFromRef2=357.04
AngleAtTarget=35.0
```

Add `-bp` (baseline bearing, Ref One to Ref Two) to switch to Full
Triangulation instead of Simple mode:

```
RangeCalc -b1 25 -b2 350 -d 250 -bp 90
```

By default the flag-driven result prints as bare `key=value` lines to
stdout. Add `-o <file>` to instead write a nicely formatted report
(inputs + results, boxed heading) to a file:

```
RangeCalc -b1 25 -b2 350 -d 250 -bp 90 -o report.txt
```

`-b1`, `-b2`, and `-bp` are compass bearings in degrees (1-360); `-d` is
the baseline distance in meters. All of `-b1`, `-b2`, and `-d` are
required together when using flags at all.
