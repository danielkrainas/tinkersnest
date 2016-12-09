package v1

import (
	"net/http"
	"regexp"

	"github.com/danielkrainas/tinkersnest/api/describe"
	"github.com/danielkrainas/tinkersnest/api/errcode"
)

var (
	IDRegex = regexp.MustCompile(`(?i)[0-9A-F]{8}-[0-9A-F]{4}-[4][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}`)

	versionHeader = describe.Parameter{
		Name:        "TinkersNest-Version",
		Type:        "string",
		Description: "The build version of the TinkersNest API server.",
		Format:      "<version>",
		Examples:    []string{"0.0.0-dev"},
	}

	hostHeader = describe.Parameter{
		Name:        "Host",
		Type:        "string",
		Description: "",
		Format:      "<hostname>",
		Examples:    []string{"api.tinkersnest.io"},
	}

	postNameParameter = describe.Parameter{
		Name:        "post_name",
		Type:        "string",
		Description: "Identifier for a post",
		Required:    true,
	}

	userNameParameter = describe.Parameter{
		Name:        "user_name",
		Type:        "string",
		Description: "Identifier for a user",
		Required:    true,
	}

	jsonContentLengthHeader = describe.Parameter{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "Length of the JSON body.",
		Format:      "<length>",
	}

	zeroContentLengthHeader = describe.Parameter{
		Name:        "Content-Length",
		Type:        "integer",
		Description: "The 'Content-Length' header must be zero and the body must be empty.",
		Format:      "0",
	}

	resourceNotFoundResp = describe.Response{
		Name:        "Resource Unknown Error",
		StatusCode:  http.StatusNotFound,
		Description: "The resource is not known to the server.",
		Headers: []describe.Parameter{
			versionHeader,
			jsonContentLengthHeader,
		},
		Body: describe.Body{
			ContentType: "application/json; charset=utf-8",
			Format:      errorsBody,
		},
		ErrorCodes: []errcode.ErrorCode{
			ErrorCodeResourceUnknown,
		},
	}
)

var (
	errorsBody = `{
	"errors:" [
	    {
            "code": <error code>,
            "message": <error message>,
            "detail": ...
        },
        ...
    ]
}`

	blogPostBody = `{
	"name": ...,
	"created": <epoch seconds>,
	"publish": true|false,
	"title": ...,
	"content": ...
}`

	blogPostListBody = `[
` + blogPostBody + `, ...
]`

	userBody = `{
	"name": ...,
	"full_name": "John Doe",
	"email": "j.doe@example.org"
}`

	userListBody = `[
` + userBody + `, ...
]`
)

var API = struct {
	Routes []describe.Route `json:"routes"`
}{
	Routes: routeDescriptors,
}

var routeDescriptors = []describe.Route{
	{
		Name:        RouteNameBase,
		Path:        "/v1",
		Entity:      "Base",
		Description: "Base V1 API route, can be used for lightweight health and version check.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Check that the server supports the TinkersNest V1 API.",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "The API implements the V1 protocol and is accessible.",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},

						Failures: []describe.Response{
							{
								Description: "The API does not support the V1 protocol.",
								StatusCode:  http.StatusNotFound,
								Headers: []describe.Parameter{
									versionHeader,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameAuth,
		Path:        "/v1/auth",
		Entity:      "JWT",
		Description: "",
		Methods: []describe.Method{
			{
				Method:      "POST",
				Description: "",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "All posts returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNamePostByName,
		Path:        "/v1/blog/posts/{post_name}",
		Entity:      "Post",
		Description: "Route to retrieve, update, and delete a single post by name.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Get a post by name",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							postNameParameter,
						},

						Successes: []describe.Response{
							{
								Description: "post returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      blogPostBody,
								},
							},
						},
					},
				},
			},
			{
				Method:      "DELETE",
				Description: "Delete a post by name",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							postNameParameter,
						},

						Successes: []describe.Response{
							{
								Description: "post removed",
								StatusCode:  http.StatusNoContent,
								Headers: []describe.Parameter{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameBlog,
		Path:        "/v1/blog/posts",
		Entity:      "[]Post",
		Description: "Route to retrieve the list of posts and create new ones.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Get all posts",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "All posts returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      blogPostListBody,
								},
							},
						},
					},
				},
			},
			{
				Method:      "POST",
				Description: "Create a blog post",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "Post created",
								StatusCode:  http.StatusCreated,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      blogPostBody,
								},
							},
						},

						Failures: []describe.Response{},
					},
				},
			},
		},
	},

	{
		Name:        RouteNameUserByName,
		Path:        "/v1/users/{user_name}",
		Entity:      "User",
		Description: "Route to retrieve, update, and delete a single user by name.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Retrieve a single user",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							userNameParameter,
						},

						Successes: []describe.Response{
							{
								Description: "user returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      userBody,
								},
							},
						},
					},
				},
			},
			{
				Method:      "PUT",
				Description: "Modify a single user",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							userNameParameter,
						},

						Successes: []describe.Response{
							{
								Description: "user returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      userBody,
								},
							},
						},
					},
				},
			},
			{
				Method:      "DELETE",
				Description: "Delete a user",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						PathParameters: []describe.Parameter{
							postNameParameter,
						},

						Successes: []describe.Response{
							{
								Description: "user deleted",
								StatusCode:  http.StatusNoContent,
								Headers: []describe.Parameter{
									versionHeader,
									zeroContentLengthHeader,
								},
							},
						},
					},
				},
			},
		},
	},
	{
		Name:        RouteNameUserRegistry,
		Path:        "/v1/users",
		Entity:      "[]User",
		Description: "Route to retrieve the list of users and create new ones.",
		Methods: []describe.Method{
			{
				Method:      "GET",
				Description: "Get all users",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "All users returned",
								StatusCode:  http.StatusOK,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      userListBody,
								},
							},
						},
					},
				},
			},
			{
				Method:      "PUT",
				Description: "Create a user",
				Requests: []describe.Request{
					{
						Headers: []describe.Parameter{
							hostHeader,
						},

						Successes: []describe.Response{
							{
								Description: "User created",
								StatusCode:  http.StatusCreated,
								Headers: []describe.Parameter{
									versionHeader,
									jsonContentLengthHeader,
								},

								Body: describe.Body{
									ContentType: "application/json; charset=utf-8",
									Format:      userBody,
								},
							},
						},

						Failures: []describe.Response{},
					},
				},
			},
		},
	},
}

var routeDescriptorsMap map[string]describe.Route

func init() {
	routeDescriptorsMap = make(map[string]describe.Route, len(routeDescriptors))
	for _, descriptor := range routeDescriptors {
		routeDescriptorsMap[descriptor.Name] = descriptor
	}
}
