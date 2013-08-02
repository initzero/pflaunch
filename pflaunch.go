package main

import (
	"syscall"
	"log"
	"unsafe"
	"os/exec"
	"os"
	"path/filepath"
	"strings"
)

import (
	"github.com/lxn/win"
    "github.com/lxn/walk"
    . "github.com/lxn/walk/declarative"
)

type DWORD uint32

type PFRegistry struct {
	EAX 						DWORD
	EregCount 					DWORD
	FastAnimations 				DWORD
	GeometryCacheSize			DWORD 	
	ListLAN						DWORD
	MultiServerListBinDerDunDat	DWORD
	PFAutoClose					DWORD
	PFClientPort				string
	PFCommandLine				string
	PFFastMulti					DWORD
	PFFastStart					DWORD
	PFServerPort				string
	ResolutionBackbufferFormat 	DWORD
	ResolutionBitDepth 			DWORD
	ResolutionHeight 			DWORD
	ResolutionWidth 			DWORD
	Vsync 						DWORD
	WindowMode 					DWORD
}

func (p *PFRegistry) Fill() (err error) {
	var h syscall.Handle
	err = syscall.RegOpenKeyEx(syscall.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr("SOFTWARE\\Volition\\Red Faction"), 0, syscall.KEY_READ, &h)
	checkErr(err, "open registry")
	defer syscall.RegCloseKey(h)

	val, err := ReadRegDWORD(h, "EAX")
	checkErr(err, "read EAX")
	p.EAX = val

	val, err = ReadRegDWORD(h, "EregCount")
	checkErr(err, "read EregCount")
	p.EregCount = val

	val, err = ReadRegDWORD(h, "Fast Animations")
	checkErr(err, "read Fast Animations")
	p.FastAnimations = val

	val, err = ReadRegDWORD(h, "Geometry Cache Size")
	checkErr(err, "read Geometry Cache Size")
	p.GeometryCacheSize = val

	val, err = ReadRegDWORD(h, "ListLAN")
	checkErr(err, "read ListLAN")
	p.ListLAN = val

	val, err = ReadRegDWORD(h, "MultiServerListBinDerDunDat")
	checkErr(err, "read MultiServerListBinDerDunDat")
	p.MultiServerListBinDerDunDat = val

	val, err = ReadRegDWORD(h, "pfAutoClose")
	checkErr(err, "read pfAutoClose")
	p.PFAutoClose = val

	val, err = ReadRegDWORD(h, "pfFastMulti")
	checkErr(err, "read pfFastMulti")
	p.PFFastMulti = val

	val, err = ReadRegDWORD(h, "pfFastStart")
	checkErr(err, "read pfFastStart")
	p.PFFastStart = val

	val, err = ReadRegDWORD(h, "Resolution Backbuffer Format")
	checkErr(err, "read Resolution Backbuffer Format")
	p.ResolutionBackbufferFormat = val

	val, err = ReadRegDWORD(h, "Resolution Bit Depth")
	checkErr(err, "read Resolution Bit Depth")
	p.ResolutionBitDepth = val

	val, err = ReadRegDWORD(h, "Resolution Height")
	checkErr(err, "read Resolution Height")
	p.ResolutionHeight = val

	val, err = ReadRegDWORD(h, "Resolution Width")
	checkErr(err, "read Resolution Width")
	p.ResolutionWidth = val

	val, err = ReadRegDWORD(h, "Vsync")
	checkErr(err, "read Vsync")
	p.Vsync = val

	val, err = ReadRegDWORD(h, "windowMode")
	checkErr(err, "read windowMode")
	p.WindowMode = val

	str, err := ReadRegString(h, "pfClientPort")
	checkErr(err, "read pfClientPort")
	p.PFClientPort = str

	str, err = ReadRegString(h, "pfCommandLine")
	checkErr(err, "read pfCommandLine")
	p.PFCommandLine = str

	str, err = ReadRegString(h, "pfServerPort")
	checkErr(err, "read pfServerPort")
	p.PFServerPort = str

	return err
}

func (p *PFRegistry) Save() (ret int32) {
	var h win.HKEY
	ret = win.RegOpenKeyEx(win.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr("SOFTWARE\\Volition\\Red Faction"), 0, syscall.KEY_WRITE, &h)
	defer win.RegCloseKey(h)

	ret = SetRegDWORD(h, "EAX", p.EAX)
	ret = SetRegDWORD(h, "EregCount", p.EregCount)
	ret = SetRegDWORD(h, "Fast Animations", p.FastAnimations)
	ret = SetRegDWORD(h, "Geometry Cache Size", p.GeometryCacheSize)
	ret = SetRegDWORD(h, "ListLAN", p.ListLAN)
	ret = SetRegDWORD(h, "MultiServerListBinDerDunDat", p.MultiServerListBinDerDunDat)
	ret = SetRegDWORD(h, "pfAutoClose", p.PFAutoClose)
	ret = SetRegDWORD(h, "pfFastMulti", p.PFFastMulti)
	ret = SetRegDWORD(h, "pfFastStart", p.PFFastStart)
	ret = SetRegDWORD(h, "Resolution Backbuffer Format", p.ResolutionBackbufferFormat)
	ret = SetRegDWORD(h, "Resolution Bit Depth", p.ResolutionBitDepth)
	ret = SetRegDWORD(h, "Resolution Height", p.ResolutionHeight)
	ret = SetRegDWORD(h, "Resolution Width", p.ResolutionWidth)
	ret = SetRegDWORD(h, "Vsync", p.Vsync)
	ret = SetRegDWORD(h, "windowMode", p.WindowMode)
	ret = SetRegString(h, "pfServerPort", p.PFServerPort)
	ret = SetRegString(h, "pfClientPort", p.PFClientPort)
	ret = SetRegString(h, "pfCommandLine", p.PFCommandLine)
	return ret
}

type PFDataSource struct {
	Registry *PFRegistry
	PFFastStartFlag bool
	PFFastMultiFlag bool
	VSyncFlag bool
	FastAnimationsFlag bool
	PFAutoCloseFlag bool
	EAXFlag bool
	PFClientPortNumber float64
	RenderingCacheNumber float64
	ServerUploadNumber float64
	SimultaneousPingsNumber float64
	ModName string
	ResolutionName string
	ColorDepthName string
	ServerTrackerName string
}

func checkErr(err error, context string) {
	if err != nil {
		log.Printf("[%s] %s\n", context, err)
	}
}

// read string from given handle for given key
func ReadRegString(h syscall.Handle, key string) (value string, err error) {
	var typ uint32
	var buf [74]uint16
	n := uint32(len(buf))
	err = syscall.RegQueryValueEx(h, syscall.StringToUTF16Ptr(key), nil, &typ, (*byte)(unsafe.Pointer(&buf[0])), &n)
	if err != nil {
		log.Println(err)
		return value, err
	}
	return syscall.UTF16ToString(buf[:]), err
}

// read dword from given handle for given key
func ReadRegDWORD(h syscall.Handle, key string) (value DWORD, err error) {
	var typ uint32
	var buf DWORD
	var n uint32 = 32
	err = syscall.RegQueryValueEx(h, syscall.StringToUTF16Ptr(key), nil, &typ, (*byte)(unsafe.Pointer(&buf)), &n)
	if err != nil {
		log.Println(err)
		return value, err
	}
	return buf, err
}

// set registry string
// Have to use the win package from lxn as golang syscall doesn't provide a RegSetValueEx function
func SetRegString(h win.HKEY, key string, value string) (ret int32) {
	buf := syscall.StringToUTF16Ptr(value)
	ret = win.RegSetValueEx(h, syscall.StringToUTF16Ptr(key), 0, win.REG_SZ, (*byte)(unsafe.Pointer(buf)), 32)
	return ret
}

func SetRegDWORD(h win.HKEY, key string, value DWORD) (ret int32) {
	ret = win.RegSetValueEx(h, syscall.StringToUTF16Ptr(key), 0, win.REG_DWORD, (*byte)(unsafe.Pointer(&value)), 4)
	return ret
}

// misc notes

// its -alttab for stretched, and either -window or -windowed for that one

// UI Models
type Model struct {
	Id int
	Name string
}

type DWORDModel struct {
	Id int
	Name string
	Value DWORD
}

type StringModel struct {
	Id int
	Name string
	Value string
}

type ColorDepthModel StringModel

func ColorDepths() []*ColorDepthModel {
	return []*ColorDepthModel{
		{0, "16bit", ""},
		{1, "32bit", ""},
	}
}

type ResolutionModel StringModel

func Resolutions() []*ResolutionModel {
	return []*ResolutionModel{
		{0, "640x480", ""},
		{1, "720x480", ""},
		{2, "720x576", ""},
		{3, "800x600", ""},
		{4, "960x600", ""},
		{5, "1024x768", ""},
		{6, "1152x864", ""},
		{7, "1280x720", ""},
		{8, "1280x768", ""},
		{9, "1280x800", ""},
		{10, "1280x960", ""},
		{11, "1280x1024", ""},
		{12, "1360x768", ""},
		{13, "1440x900", ""},
		{14, "1440x1080", ""},
		{15, "1600x1200", ""},
		{16, "1680x1050", ""},
		{17, "1920x1080", ""},
		{18, "1920x1200", ""},
	}
}

type DisplayMode DWORDModel

func DisplayModes() []*DisplayMode {
	return []*DisplayMode{
		{0, "Normal: Default RF Fullscreen Mode", DWORD(0x01)},
		{1, "Windowed: Draggable fixed size window", DWORD(0x02)},
		{2, "Stretched: Window is stretched to your desktop resolution", DWORD(0x01)},
	}
}


// RF mods will be in the Red Faction/mods/<modname> folder
// Launch with mod using PF.exe -mod <modname>
type RFMod Model

// Walk "mods" directory, pulling out only the first level
func RFMods() []*RFMod {
	mods := []*RFMod{{0, "None"}}
	skip := make([]string, 10)
	index := 1
	err := filepath.Walk("mods", func(path string, info os.FileInfo, errIn error) (err error) {
			if info.IsDir() && path != "mods" {
				for _, v := range skip {
					if v != "" && strings.Contains(path, v) {
						return err
					}
				}
				mods = append(mods, &RFMod{index, info.Name()})
				index += 1
				skip = append(skip, path)
			}
			return err
		})
	checkErr(err, "filepath walk")
	return mods
}

func main() {
	ds := &PFDataSource{}
	registry := &PFRegistry{}
	registry.Fill()
	ds.Registry = registry
	ds.PFClientPortNumber = float64(7755)
	ds.RenderingCacheNumber = float64(8)
	ds.ServerUploadNumber = float64(1000000)
	ds.SimultaneousPingsNumber = float64(3)
	ds.ServerTrackerName = "thq.multiplay.net"
	ds.ResolutionName = "640x480"
	ds.ColorDepthName = "32bit"
	log.Printf("%#v", registry)

	mw, err := walk.NewMainWindow()
	checkErr(err, "new window")
	icon, err := walk.NewIconFromFile("PureIcon.ico")
	checkErr(err, "load icon")

	var db *walk.DataBinder

   	dMw := MainWindow{
   		AssignTo: &mw,
        Title:   "Pure Faction",
        MinSize: Size{640, 500},
        Layout:  VBox{},
		DataBinder: DataBinder{
			DataSource: ds,
			AssignTo: &db,
		},
        Children: []Widget{
        	TabWidget{
        		Pages: []TabPage{
        			TabPage{
        				Title:"Server Browser", 
        				Layout: VBox{MarginsZero: true}, 
        				Children: []Widget{
        					WebView{URL: "http://rm.servehttp.com/viewpage.php"},
        				},
        			},
        			TabPage{
        				Title:"Multiplayer", 
        				Layout: VBox{MarginsZero: true},
        				Children: []Widget{
        					Composite{
        						Layout: Grid{},
        						Children: []Widget{
		        					GroupBox{
		        						Title: "Select window mode:",
		        						Layout: Grid{},
		        						Column: 0,
		        						ColumnSpan: 1,
		        						Row: 0,
		        						Children: []Widget{
		        							ComboBox{
		        								Model: DisplayModes(),
		        								BindingMember: "Value",
		        								DisplayMember: "Name",
		        								CurrentIndex: 0,
		        								Value: Bind("Registry.WindowMode", SelRequired{}),
		        								OnCurrentIndexChanged: func() { db.Submit(); log.Printf("%#v", ds.Registry) },
		        							},
		        						},
		        					},
		        					GroupBox{
		        						Title: "Behavior:",
		        						Layout: Grid{},
		        						Column: 0,
		        						Row: 1,
		        						Children: []Widget{
		        							CheckBox{
		        								Text: "Fast Start: Launch the game faster (skips intros and loading delays)",
		        								Column: 0,
		        								Row: 0,	
		        								Checked: Bind("PFFastStartFlag"),
		        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
		        							},
		        							CheckBox{
		        								Text: "Fast Multi: Go straight to the server list when game starts",
		        								Column: 0,
		        								Row: 1,
		        								Checked: Bind("PFFastMultiFlag"),
		        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
		        							},
		        							CheckBox{
		        								Text: "Auto Close: Close launcher when starting the game",
		        								Column: 0,
		        								Row: 2,
		        								Checked: Bind("PFAutoCloseFlag"),
		        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
		        							},
		        						},
		        					},
		        					GroupBox{
		        						Title: "Misc Settings:",
		        						Layout: Grid{},
		        						Column: 0,
		        						Row: 2,
		        						Children: []Widget{
		        							Label{
		        								Column: 0,
		        								ColumnSpan: 1,
		        								Row: 0,
		        								Text: "Client Port:",
		        							},
		        							NumberEdit{
		        								Column: 1,
		        								Row: 0,
		        								ColumnSpan: 1,
		        								MinValue: 1024,
		        								MaxValue: 65535,
		        								Value: Bind("PFClientPortNumber"),
		        								OnValueChanged: func() { db.Submit(); log.Printf("%#v", ds.Registry) },
		        							},
  	      									Label{
        										Text: "Tracker: ",
        										Column: 0,
        										Row: 1,
        									},
        									LineEdit{
        										Column: 1,
        										Row: 1,
        										Text: Bind("ServerTrackerName"),
        										OnTextChanged: func() { db.Submit(); log.Printf("%#v", ds) },
        										MaxLength: 1024,
        									},
		        							Label{
		        								Column: 0,
		        								ColumnSpan: 1,
		        								Row: 2,
		        								Text: "Extra commands:",
		        							},
		        							TextEdit{
		        								Column: 1,
		        								ColumnSpan: 1,
		        								Row: 2,
		        								Text: Bind("Registry.PFCommandLine"),
		        							},
		        						},
		        					},
		        					GroupBox{
		        						Title: "Mod:",
		        						Layout: Grid{},
		        						Column: 0,
		        						Row: 3,
		        						Children: []Widget{
					        				ComboBox{
				        						Column: 1,
				        						Row: 3,
				        						Model: RFMods(),
				        						DisplayMember: "Name",
				        						BindingMember: "Name",
				        						Value: Bind("ModName"),
				        						CurrentIndex: 0,
				        					},
				        				},
		        					},
		        					Composite{
		        						Layout: HBox{},
		        						Column: 1,
		        						Row: 3,
		        						Children: []Widget{
		        							PushButton{
		        								Text: "Save",
		        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds.Registry) },
		        							},
		        							PushButton{
		        								Text: "Launch",
		        								OnClicked: func() {
		        									go func() {
		        										err := exec.Command("PF.exe").Run()
		        										checkErr(err, "Launch")
		        									}()
		        								},
		        							},
		        						},
		        					},
        						},
        					},
        				},
        			},
        			TabPage{
        				Title:"Game Settings", 
        				Layout: VBox{MarginsZero: true},
        				Children: []Widget{
        					Composite{
        						Layout: Grid{},
        						Children: []Widget{
        							GroupBox{
        								Title: "Resolution / Color Depth:",
        								Column: 0,
        								Row: 0,
        								Layout: Grid{},
        								Children: []Widget{
        									ComboBox{
		        								Column: 0,
		        								Row: 0,
		        								Model: Resolutions(),
		        								DisplayMember: "Name",
		        								BindingMember: "Name",
		        								Value: Bind("ResolutionName"),
		        								CurrentIndex: 0,
		        							},
        									ComboBox{
		        								Column: 1,
		        								Row: 0,
		        								Model: ColorDepths(),
		        								DisplayMember: "Name",
		        								BindingMember: "Name",
		        								Value: Bind("ColorDepthName"),
		        								CurrentIndex: 0,
		        							},
        								},
        							},
        							GroupBox{
        								Title: "Video Options",
        								Column: 0,
        								Row: 1,
        								Layout: Grid{},
        								Children: []Widget{
        									Composite{
        										Layout: HBox{},
        										Column: 0,
        										Row: 0,
        										Children: []Widget{
				        							Label{
				        								Text: "Rendering Cache: ",
				        							},
				        							NumberEdit{
				        								MinValue: 2,
				        								MaxValue: 32,
				        								MinSize: Size{75,100},
				        								Suffix: " (MB)",
				        								Value: Bind("RenderingCacheNumber"),
				        								OnValueChanged: func() { db.Submit(); log.Printf("%#v", ds) },
				        							},
				        							CheckBox{
				        								Text: "VSync",
				        								Checked: Bind("VSyncFlag"),
				        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
				        							},
				        							CheckBox{
				        								Text: "Fast Animations",
				        								Checked: Bind("FastAnimationsFlag"),
				        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
				        							},
				        							HSpacer{

				        							},
        										},
        									},
        								},
        							},
        							GroupBox{
        								Title: "Network Options:",
        								Column: 0,
        								Row: 2,
        								Layout: Grid{},
        								Children: []Widget{
        									Composite{
        										Column: 0,
        										Row: 0,
        										Layout: HBox{},
        										Children: []Widget{
    												Label{
        												Text: "Server Upload Limit: ",
		        									},
		        									NumberEdit{
		        										MinValue: 1,
		        										MaxValue: 1000000,
		        										MinSize: Size{100,50},
		        										Suffix: " (B/s)",
		        										Value: Bind("ServerUploadNumber"),
		        										OnValueChanged: func() { db.Submit(); log.Printf("%#v", ds) },
		        									},
		        									Label{
		        										Text: "Simultaneous Pings: ",
		        									},
		        									NumberEdit{
		        										MinValue: 3,
		        										MaxValue: 1024,
		        										Value: Bind("SimultaneousPingsNumber"),
		        										OnValueChanged: func() { db.Submit(); log.Printf("%#v", ds) },
		        									},
		        									HSpacer{
		        									},
    											},
        									},
        								},
        							},
        							GroupBox{
        								Title: "Misc Options:",
        								Column: 0,
        								Row: 3,
        								Layout: Grid{},
        								Children: []Widget{
        									Composite{
        										Column: 0,
        										Row: 0,
        										Layout: VBox{},
        										Children: []Widget{
		        									CheckBox{
				        								Text: "Enable EAX Sound",
				        								Checked: Bind("EAXFlag"),
				        								OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
		        									},
        											PushButton{
        												Text: "Open Maps Folder",
        												OnClicked: func() {
        													go func() {
        														mapsFolder, err := filepath.Abs("user_maps/multi")
        														log.Println(mapsFolder)
        														checkErr(err, "find mapsFolder")
        														err = exec.Command("explorer", mapsFolder).Run()
        														checkErr(err, "open maps")
        													}()
        												},
        											},
    											},
        									},
        								},
        							},
        							Composite{
        								Layout: HBox{},
        								Column: 1,
        								Row: 3,
        								Children: []Widget{
        									PushButton{
        										Text: "Save",
        										OnClicked: func() { db.Submit(); log.Printf("%#v", ds) },
        									},
        								},
        							},
        						},
        					},
        				},
        			},
        			TabPage{
        				Title:"Server Setup", 
        				Layout: VBox{MarginsZero: true},
        			},
        			TabPage{
        				Title:"About", 
        				Layout: VBox{MarginsZero: true},
        			},
        		},
        	},
        },
    }
 
    // weird workaround to set the window icon, as you can't do it with declarative for some reason?
    go func() {
    	for {
    		if mw.Icon() == nil {
    			mw.SetIcon(icon)
    		} else {
    			break
    		}
    	}
    }()

    _, runErr := dMw.Run()
    if runErr != nil {
    	log.Println(runErr)
    }
}