#ShaleAppTest

A sample application to demonstrate Golang aptitude. The app uses and channels and websockets to provide a real-time experience.

##To run

###Prepare

```bash
make install-deps #attempts to install Glide
make bootstrap #installs deps from glide.yaml
```

If for some reason the `make install-deps` command fails you can find glide on github and install it manually or simply run `go get github.com/gorilla/websocket` and copy the dependency into a vendor directory at the root of this project.

###Build

```bash
make build
```

###Run

```bash
make run
```

The output of this command will show you what URL to use in your browser to see the app. In order to see the messages sync with multiple users you should now open up another browser or an incognito window and visit the application there as well.