# Octops Image Syncer

Watches deployments of Agones Fleets and pre-pull gameserver images onto every Kubernetes worker node.

## Use Cases
- Pre-Pull images on every Kubernetes worker node before a gameserver replica is schedule
- Lower startup time of gameservers on nodes that never had a replica of the gameserver scheduled on it
- Fleets which are scheduled up and down using an on demand or custom scheduling strategy

## How it works

The Octops Image Syncer will be notified by the Kubernetes cluster everytime a Fleet is deployed, updated or deleted. This is achieved by using the https://github.com/Octops/agones-event-broadcaster and the broker interface.

When an event is received, the syncer ([Daemonset](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)) will check if the Fleet `image:tag` is present on the node and if not it will request the image to be pulled by the container runtime. 

The syncing process is not coupled to any specific runtime like Docker. The Octops Image Syncer works on any Kubernetes cluster that supports the [Container Runtime Interface](https://kubernetes.io/docs/setup/production-environment/container-runtimes/).

A few extra considerations:
- It runs as Daemonset and pulls images for all fleets deployed on the cluster. Labels or annotations for filtering are not supported yet.
- The Dockerfile used to build the application uses a [Distroless image](https://github.com/GoogleContainerTools/distroless). Distroless images contain only your application and its runtime dependencies. They do not contain package managers, shells or any other programs you would expect to find in a standard Linux distribution.
- Tested on Kubernetes 1.20+

Output:

```bash
#Image not present
{"fleet":"simple-game-server","image":"gcr.io/agones-images/simple-game-server:0.3","message":"fleet synced","ref":"sha256:f8cdc89145cb0b5d6ee2ea95968310c45e4f453dd24ac682ff13f50f0d4b921d","severity":"info","time":"2021-09-18T10:42:12.456488899Z"}

#Image already present
{"fleet":"simple-game-server","image":"gcr.io/agones-images/simple-game-server:0.3","message":"image already present","severity":"info","time":"2021-09-18T10:41:56.179672762Z"}
```

## Requirements

- Kubernetes cluster that supports the [Container Runtime Interface - CRI](https://kubernetes.io/blog/2016/12/container-runtime-interface-cri-in-kubernetes/)
- Recommended 1.20+
- Access to the CRI service using Unix Domain Sockets or gRPC

## Security Concerns
- Host path and sockets mounts. Make sure your company security policies allow you to mount volumes pointing to host path. That is usually a bad idea and must be evaluated carefully.
- The Octops Image Syncer is still a working in progress. As any other software is not free of possibly supply chain attack or any package that might me compromised.

## Build and Install
- Build your own docker image and push to your registry
- Deploy the Daemonset on each cluster you want the application watching Fleets and syncing images
- Example of manifest present on hack/install.yaml. Update the image to reflect your repo, image name and tag

### Container Runtime Interface, Mounts and Volumes
The Octops Image Syncer has to establish a connection with the container runtime layer to perform the image check and image pull operations. Depending on your setup, cloud provider, kubernetes version and security policies the communication method or socket path can vary. 

Make sure that the manifest that deploys the Octops Image Syncer reflects your environment and container runtime settings. 

Set the `path` of the volume to the socket that the container runtime will accept connections.

```yaml
volumes:
- name: runtime-sock
  hostPath:
    path: "/run/containerd/containerd.sock"
    type: Socket
```

The information about the ENV and VolumeMount don't have to be changed. This configuration will work for any socket path.

```yaml
# Don't change this information unless you know exactly what you are doing. On very rare situations, the container runtime exposes a gRPC endpoint instead of a Unix Socket. 
# This is not standard and requires settings to be manually changed by an operator. Not recommended for production environments
env:
- name: CONN_TARGET
  value: "unix:///run/runtime/cri.sock"
volumeMounts:
- mountPath: /run/runtime/cri.sock
  name: runtime-sock
  readOnly: false
```

Below is listed 3 possible container runtimes and its Unix Socket's paths:

- Containerd: /run/containerd/containerd.sock

- Docker: /var/run/docker.sock

- CRI-O: /var/run/crio/crio.sock

### Extra
If you are running your cluster using [k3s](https://k3s.io/) the containerd socket is located on `/run/k3s/containerd/containerd.sock`.

Be aware the support for Docker runtime will be deprecated. You can find more information on https://kubernetes.io/blog/2020/12/02/dont-panic-kubernetes-and-docker/