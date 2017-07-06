#ShaleAppTest

A sample application to demonstrate Golang aptitude.

##To run

###Prepare

```bash
make install-deps #attempts to install Glide
make bootstrap #installs deps from glide.yaml
```

###Build

```bash
make build
```

###Run

```bash
make run
```

The output of this command will show you what URL to use in your browser to see the app. In order to see the messages sync with multiple users you should now open up another browser or an incognito window and visit the application there as well.