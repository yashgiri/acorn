---
title: "acorn run"
---
## acorn run

Run an app from an image or Acornfile

```
acorn run [flags] IMAGE|DIRECTORY [acorn args]
```

### Examples

```
# Publish and Expose Port Syntax
  # Publish port 80 for any containers that define it as a port
  acorn run -p 80 .

  # Publish container "myapp" using the hostname app.example.com
  acorn run --publish app.example.com:myapp .

  # Expose port 80 to the rest of the cluster as port 8080
  acorn run --expose 8080:80/http .

# Labels and Annotations Syntax
  # Add a label to all resources created by the app
  acorn run --label key=value .

  # Add a label to resources created for all containers
  acorn run --label containers:key=value .

  # Add a label to the resources created for the volume named "myvolume"
  acorn run --label volumes:myvolume:key=value .

# Link Syntax
  # Link the running acorn application named "mydatabase" into the current app, replacing the container named "db"
  acorn run --link mydatabase:db .

# Secret Syntax
  # Bind the acorn secret named "mycredentials" into the current app, replacing the secret named "creds". See "acorn secrets --help" for more info
  acorn run --secret mycredentials:creds .

# Volume Syntax
  # Create the volume named "mydata" with a size of 5 gigabyes and using the "fast" storage class
  acorn run --volume mydata,size=5G,class=fast .

  # Bind the acorn volume named "mydata" into the current app, replacing the volume named "data", See "acorn volumes --help for more info"
  acorn run --volume mydata:data .

# Automatic upgrades
  # Automatic upgrade for an app will be enabled if '#', '*', or '**' appears in the image's tag. Tags will sorted according to the rules for these special characters described below. The newest tag will be selected for upgrade.

  # '#' denotes a segment of the image tag that should be sorted numerically when finding the newest tag.
  # This example deploys the hello-world app with auto-upgrade enabled and matching all major, minor, and patch versions:
  acorn run myorg/hello-world:v#.#.#

  # '*' denotes a segment of the image tag that should sorted alphabetically when finding the latest tag.
  # In this example, if you had a tag named alpha and a tag named zeta, zeta would be recognized as the newest:
  acorn run myorg/hello-world:*

  # '**' denotes a wildcard. This segment of the image tag won't be considered when sorting. This is useful if your tags have a segment that is unpredictable.
  # This example would sort numerically according to major and minor version (ie v1.2) and ignore anything following the "-":
  acorn run myorg/hello-world:v#.#-**

  # NOTE: Depending on your shell, you may see errors when using '*' and '**'. Using quotes will tell the shell to ignore them so Acorn can parse them:
  acorn run "myorg/hello-world:v#.#-**"

  # Automatic upgrades can be configured explicitly via a flag.
  # In this example, the tag will always be "latest", but acorn will periodically check to see if new content has been pushed to that tag:
  acorn run --auto-upgrade myorg/hello-world:latest

  # To have acorn notify you that an app has an upgrade available and require confirmation before proceeding, set the notify-upgrade flag:
  acorn run --notify-upgrade myorg/hello-world:v#.#.# myapp
  # To proceed with an upgrade you've been notified of:
  acorn update --confirm-upgrade myapp

```

### Options

```
      --annotation strings        Add annotations to the app and the resources it creates (format [type:][name:]key=value) (ex k=v, containers:k=v)
      --auto-upgrade              Enabled automatic upgrades.
  -b, --bidirectional-sync        In interactive mode download changes in addition to uploading
      --compute-class strings     Set computeclass for a workload in the format of workload=computeclass. Specify a single computeclass to set all workloads. (ex foo=example-class or example-class)
  -i, --dev                       Enable interactive dev mode: build image, stream logs/status in the foreground and stop on exit
  -e, --env strings               Environment variables to set on running containers
  -f, --file string               Name of the build file (default "DIRECTORY/Acornfile")
  -h, --help                      help for run
      --interval string           If configured for auto-upgrade, this is the time interval at which to check for new releases (ex: 1h, 5m)
  -l, --label strings             Add labels to the app and the resources it creates (format [type:][name:]key=value) (ex k=v, containers:k=v)
      --link strings              Link external app as a service in the current app (format app-name:container-name)
  -m, --memory strings            Set memory for a workload in the format of workload=memory. Only specify an amount to set all workloads. (ex foo=512Mi or 512Mi)
  -n, --name string               Name of app to create
      --notify-upgrade            If true and the app is configured for auto-upgrades, you will be notified in the CLI when an upgrade is available and must confirm it
  -o, --output string             Output API request without creating app (json, yaml)
      --profile strings           Profile to assign default values
  -p, --publish strings           Publish port of application (format [public:]private) (ex 81:80)
  -P, --publish-all               Publish all (true) or none (false) of the defined ports of application
  -q, --quiet                     Do not print status
      --region string             Region in which to deploy the app, immutable
  -s, --secret strings            Bind an existing secret (format existing:sec-name) (ex: sec-name:app-secret)
      --target-namespace string   The name of the namespace to be created and deleted for the application resources
  -u, --update                    Update the app if it already exists
  -v, --volume stringArray        Bind an existing volume (format existing:vol-name,field=value) (ex: pvc-name:app-data)
      --wait                      Wait for app to become ready before command exiting (default true)
```

### Options inherited from parent commands

```
  -A, --all-projects        Use all known projects
      --debug               Enable debug logging
      --debug-level int     Debug log level (valid 0-9) (default 7)
      --kubeconfig string   Explicitly use kubeconfig file, overriding current project
  -j, --project string      Project to work in
```

### SEE ALSO

* [acorn](acorn.md)	 - 

