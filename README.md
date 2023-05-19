
API:
https://e94b1ce1-c456-4c23-ad7b-863c8184c496.zeus.fyi/health

https://medium.com/@zeusfyi/zeus-ui-no-code-kubernetes-authenticated-api-tutorial-c468d5ef0446

You can create your own API following this tutorial, replace the docker image to this or your own

DockerImage: zeusfyi/blobme:latest
Args: -c, main
I also removed the auth url, so the API is public, you can add it back if you want to make it authenticated.