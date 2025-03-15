# flowfuse-agent
An agent to run FlowFuse managed instances of Node-RED microservices on PoT

```.yaml
apiVersion: datasance.com/v3
kind: Application
metadata:
  name: flowfuse-agent
spec:
  microservices:
    - name: flowfuse-agent-1
      agent:
        name: demo-1
      images:
        registry: 1
        catalogItemId: null
        x86: emirhandurmus/flowfuse-agent:3.1.3
        arm: emirhandurmus/flowfuse-agent:3.1.3
      container:
        rootHostAccess: false
        cdiDevices: []
        ports:
          - internal: 1880
            external: 1880
            protocol: tcp
        volumes:
          - hostDestination: flowfuse-agent
            containerDestination: /opt/flowfuse-device
            accessMode: rw
            type: volume
        env: []
        extraHosts: []
        commands: []
      config:
        deviceId: 
        token: 
        credentialSecret: 
        forgeURL: 
        brokerURL:
        brokerUsername:
        brokerPassword:

```
