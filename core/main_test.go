package autospotting

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func Test_getRegions(t *testing.T) {

	tests := []struct {
		name    string
		ec2conn mockEC2
		want    []string
		wantErr error
	}{{
		name: "return some regions",
		ec2conn: mockEC2{
			dro: &ec2.DescribeRegionsOutput{
				Regions: []*ec2.Region{
					{RegionName: aws.String("foo")},
					{RegionName: aws.String("bar")},
				},
			},
			drerr: nil,
		},
		want:    []string{"foo", "bar"},
		wantErr: nil,
	},
		{
			name: "return an error",
			ec2conn: mockEC2{
				dro: &ec2.DescribeRegionsOutput{
					Regions: []*ec2.Region{
						{RegionName: aws.String("foo")},
						{RegionName: aws.String("bar")},
					},
				},
				drerr: fmt.Errorf("fooErr"),
			},
			want:    nil,
			wantErr: fmt.Errorf("fooErr"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, err := getRegions(tt.ec2conn)
			CheckErrors(t, err, tt.wantErr)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRegions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_spotEnabledIsAddedByDefault(t *testing.T) {

	tests := []struct {
		name   string
		config Config
		want   string
	}{
		{
			name:   "Default No ASG Tags",
			config: Config{},
			want:   "spot-enabled=true",
		},
		{
			name: "Specified ASG Tags",
			config: Config{
				FilterByTags: "environment=dev",
			},
			want: "environment=dev",
		},
		{
			name: "Specified ASG that is just whitespace",
			config: Config{
				FilterByTags: "         ",
			},
			want: "spot-enabled=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			addDefaultFilter(&tt.config)

			if !reflect.DeepEqual(tt.config.FilterByTags, tt.want) {
				t.Errorf("addDefaultFilter() = %v, want %v", tt.config.FilterByTags, tt.want)
			}
		})
	}
}
