package cmd

import (
	"github.com/camerondurham/ch/cmd/util"
	"github.com/camerondurham/ch/cmd/util/mocks"
	"github.com/golang/mock/gomock"
	"testing"
)

const (
	validWindows1    = "C:\\Users\\csci104\\:/"
	validWindowsSub1 = "C:\\Users\\csci104\\"

	invalidWindows1    = "C\\Users\\csci104:/"
	invalidWindowsSub1 = "C\\Users\\csci104"

	validWindowsAbsSub2 = "C:\\Users\\currentDir"
	validWindows2       = "currentDir:/"
	validWindowsSub2    = "currentDir"
	validWindowsExpect2 = "C:\\Users\\currentDir:/"

	invalidWindows2    = "Users\\csci104:/"
	invalidWindowsSub2 = "Users\\csci104"

	invalidWindows3    = "C:\\Users\\csci104/"
	invalidWindowsSub3 = "C"

	invalidWindows4    = "C\\Users\\csci104"
	invalidWindowsSub4 = "C\\Users\\csci104"

	invalidWindows5 = "C\\Users\\csci104/root/projects"

	validUnix1    = "/Users/user:/"
	validUnixSub1 = "/Users/user"

	invalidUnix1    = "Users/user:/"
	invalidUnixSub1 = "Users/user"

	validUnix2       = "/Users/user/path/..:/"
	validUnixExpect2 = "/Users/user:/"
	validUnixSub2    = "/Users/user/path/.."
	validUnixAbsSub2 = "/Users/user"
)

func Test_parseHostContainerPath(t *testing.T) {
	type args struct {
		pathStr string
		v       util.Validate
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewMockValidate(ctrl)

	m.EXPECT().ValidPath(gomock.Eq(validWindowsSub1)).Return(true).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(validWindowsSub1)).Return(validWindowsSub1).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(invalidWindowsSub1)).Return(false).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(invalidWindowsSub1)).Return(invalidWindowsSub1).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(validWindowsAbsSub2)).Return(true).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(validWindowsSub2)).Return(validWindowsAbsSub2).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(invalidWindowsSub2)).Return(false).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(invalidWindowsSub2)).Return(invalidWindowsSub2).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(invalidWindowsSub3)).Return(false).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(invalidWindowsSub3)).Return(invalidWindowsSub3).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(invalidWindowsSub4)).Return(false).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(invalidWindowsSub4)).Return(invalidWindowsSub4).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(validUnixSub1)).Return(true).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(validUnixSub1)).Return(validUnixSub1).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(invalidUnixSub1)).Return(false).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(invalidUnixSub1)).Return(invalidUnixSub1).AnyTimes()

	m.EXPECT().ValidPath(gomock.Eq(validUnixSub2)).Return(true).AnyTimes()
	m.EXPECT().GetAbs(gomock.Eq(validUnixSub2)).Return(validUnixAbsSub2).AnyTimes()

	tests := []struct {
		name                     string
		args                     args
		wantHostContainerAbsPath string
		wantErr                  bool
	}{
		{
			"valid abs path windows",
			args{
				pathStr: validWindows1,
				v:       m,
			},
			validWindows1,
			false,
		},
		{
			"windows: invalid drive syntax",
			args{
				pathStr: invalidWindows1,
				v:       m,
			},
			"",
			true,
		},
		{
			"valid abs path not provided",
			args{
				pathStr: validWindows2,
				v:       m,
			},
			validWindowsExpect2,
			false,
		},
		{
			"windows: no drive provided",
			args{
				pathStr: invalidWindows2,
				v:       m,
			},
			"",
			true,
		},
		{
			"no container path",
			args{
				pathStr: invalidWindows3,
				v:       m,
			},
			"",
			true,
		},
		{
			// TODO: possibly remove, this branch is already covered
			"windows invalid drive and container path",
			args{
				pathStr: invalidWindows4,
				v:       m,
			},
			"",
			true,
		},
		{
			// TODO: possibly remove, this branch is already covered
			"windows invalid drive and container path",
			args{
				pathStr: invalidWindows5,
				v:       m,
			},
			"",
			true,
		},
		{
			"proper unix path on macOS",
			args{
				pathStr: validUnix1,
				v:       m,
			},
			validUnix1,
			false,
		},
		{
			"invalid relative path",
			args{
				pathStr: invalidUnix1,
				v:       m,
			},
			"",
			true,
		},
		{
			"valid relative unix path",
			args{
				pathStr: validUnix2,
				v:       m,
			},
			validUnixExpect2,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHostContainerAbsPath, err := parseHostContainerPath(tt.args.pathStr, tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseHostContainerPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotHostContainerAbsPath != tt.wantHostContainerAbsPath {
				t.Errorf("parseHostContainerPath() gotHostContainerAbsPath = %v, want %v", gotHostContainerAbsPath, tt.wantHostContainerAbsPath)
			}
		})
	}
}
