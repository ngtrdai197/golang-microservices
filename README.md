Create a root directory for your project.
Inside the root directory, create these subdirectories: cmd, pkg, api, internal, scripts.
1. cmd directory: This is where the application is bootstrapped. It contains the application's entry points.
2. pkg directory: This is where the public libraries live. They are the code that can be used by other applications.
3. api directory: This is where the API definitions and protocol files (like .proto, .graphql files) live.
4. internal directory: This is where the private application and library code lives. It's your actual application.
5. scripts directory: This is where the scripts to perform various build, install, analysis, etc operations live.

/myapp
  /cmd
    /myapp
      main.go
  /pkg
    /mypkg
      mypkg.go
      mypkg_test.go
  /api
  /internal
    /app
      /myapp
        /server
          server.go
        /client
          client.go
    /pkg
      /mypkg
        mypkg.go
        mypkg_test.go
  /scripts