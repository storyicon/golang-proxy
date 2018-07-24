package business

import (
	"testing"

	"github.com/storyicon/golang-proxy/model"
)

func Test_getProxyString(t *testing.T) {
	type args struct {
		p *model.Proxy
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test0",
			args: args{
				p: &model.Proxy{
					IP:     "127.0.0.1",
					Port:   "8888",
					Scheme: "http",
				},
			},
			want: "http://127.0.0.1:8888",
		},
		{
			name: "test1",
			args: args{
				p: &model.Proxy{
					IP:     "127.0.0.1:8888",
					Port:   "8888",
					Scheme: "http",
				},
			},
			want: "http://127.0.0.1:8888",
		},
		{
			name: "test2",
			args: args{
				p: &model.Proxy{
					IP:     "http://127.0.0.1",
					Port:   "8888",
					Scheme: "http",
				},
			},
			want: "http://127.0.0.1:8888",
		},
		{
			name: "test3",
			args: args{
				p: &model.Proxy{
					IP:     "http://127.0.0.1",
					Port:   "8888",
					Scheme: "",
				},
			},
			want: "http://127.0.0.1:8888",
		},
		{
			name: "test4",
			args: args{
				p: &model.Proxy{
					IP:     "",
					Port:   "8888",
					Scheme: "http",
				},
			},
			want: "",
		},
		{
			name: "test5",
			args: args{
				p: &model.Proxy{
					IP:     "1",
					Port:   "8888",
					Scheme: "http",
				},
			},
			want: "",
		},
		{
			name: "test6",
			args: args{
				p: &model.Proxy{
					IP:     "127.0.0.1:8888",
					Port:   "",
					Scheme: "",
				},
			},
			want: "http://127.0.0.1:8888",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProxyString(tt.args.p); got != tt.want {
				t.Errorf("getProxyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTemplateRender(t *testing.T) {
	type args struct {
		template string
		key      string
		value    interface{}
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test0",
			args: args{
				template: "http://baidu.com?p=${page}&s=1",
				key:      "page",
				value:    "1",
			},
			want: "http://baidu.com?p=1&s=1",
		},
		{
			name: "test1",
			args: args{
				template: "http://baidu.com?p=${page}&s=1",
				key:      "page",
				value:    1,
			},
			want: "http://baidu.com?p=1&s=1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TemplateRender(tt.args.template, tt.args.key, tt.args.value); got != tt.want {
				t.Errorf("TemplateRender() = %v, want %v", got, tt.want)
			}
		})
	}
}
