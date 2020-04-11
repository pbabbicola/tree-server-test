package tree

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"text/template"
)

func Test_getFavoriteTree(t *testing.T) {
	tests := []struct {
		name string
		args *http.Request
		want string
	}{
		{
			name: "0 favorite trees",
			args: httptest.NewRequest(http.MethodGet, "http://localhost?favoriteTree=", nil),
			want: "",
		},
		{
			name: "not even a query param",
			args: httptest.NewRequest(http.MethodGet, "http://localhost", nil),
			want: "",
		},
		{
			name: "1 favorite tree",
			args: httptest.NewRequest(http.MethodGet, "http://localhost?favoriteTree=magnolia", nil),
			want: "magnolia",
		},
		{
			name: "2 favorite trees",
			args: httptest.NewRequest(http.MethodGet, "http://localhost?favoriteTree=maple&favoriteTree=pine", nil),
			want: "maple",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := getFavoriteTree(tt.args); got != tt.want {
				t.Errorf("getFavoriteTree() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_formatResponse(t *testing.T) {
	tests := []struct {
		name string
		tree string
		want string
	}{
		{
			name: "empty",
			tree: "",
			want: "Please tell me your favorite tree.",
		},
		{
			name: "not empty",
			tree: "magnolia",
			want: "It's nice to know that your favorite tree is a magnolia.",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := formatResponse(tt.tree); got != tt.want {
				t.Errorf("formatResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_TreeHandler_renderResponse(t *testing.T) {
	tmpl := template.Must(template.New("index").Parse(`
	 <html>
		<head>
			<title>Tree Example</title>
			<meta charset="utf-8">
		</head>
		<body>
			<main>
				<p>{{ .Text }}</p>
			</main>
		</body>
	 </html>
	 `))
	handler := NewHandler(tmpl)

	type args struct {
		w    *httptest.ResponseRecorder
		text string
	}

	tests := []struct {
		name string
		args args
		code int
	}{
		{
			name: "empty",
			args: args{
				text: "Please tell me your favorite tree.",
				w:    httptest.NewRecorder(),
			},
			code: 200,
		},
		{
			name: "not empty",
			args: args{
				text: "It's nice to know that your favorite tree is a magnolia.",
				w:    httptest.NewRecorder(),
			},
			code: 200,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			handler.renderResponse(tt.args.w, tt.args.text)

			//nolint:bodyclose No need to close status httptest.ResponseRecorder result
			if tt.code != tt.args.w.Result().StatusCode {
				t.Fatalf("expected code %v but got code %v", tt.code, tt.args.w.Result().StatusCode)
			}
		})
	}
}

func Test_TreeHandler_ServeHTTP(t *testing.T) {
	tmpl := template.Must(template.New("index").Parse(`
	 <html>
		<head>
			<title>Tree Example</title>
			<meta charset="utf-8">
		</head>
		<body>
			<main>
				<p>{{ .Text }}</p>
			</main>
		</body>
	 </html>
	 `))
	handler := NewHandler(tmpl)

	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}

	tests := []struct {
		name string
		args args
		code int
	}{
		{
			name: "correct",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "http://localhost/?favoriteTree=magnolia", nil),
				w: httptest.NewRecorder(),
			},
			code: http.StatusOK,
		},
		{
			name: "wrong method",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "http://localhost/?favoriteTree=magnolia", nil),
				w: httptest.NewRecorder(),
			},
			code: http.StatusMethodNotAllowed,
		},
		{
			name: "wrong address",
			args: args{
				r: httptest.NewRequest(http.MethodGet, "http://localhost/tree/?favoriteTree=magnolia", nil),
				w: httptest.NewRecorder(),
			},
			code: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			handler.ServeHTTP(tt.args.w, tt.args.r)

			//nolint:bodyclose No need to close status httptest.ResponseRecorder result
			if tt.code != tt.args.w.Result().StatusCode {
				t.Fatalf("expected code %v but got code %v", tt.code, tt.args.w.Result().StatusCode)
			}
		})
	}
}
