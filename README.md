# Drone Google Cloud Run

A plugin to [Drone](https://drone.io) plugin to do  with [Google Cloud Run](https://cloud.google.com/run/) services.

__IMPORTANT__: This plugin currently supports only Managed Google Cloud Run services.

## Usage

The following settings changes this plugin's behavior.

* `service_account_json`: The Google Cloud Service Account JSON that has required permissions to create, update and delete Google Cloud Run services .
* `project`: The Google project where the Google Cloud Run service will be deployed.
* `region`: The Google Cloud region e.g asia-south1 where the Google Cloud Run service will be deployed.
* `service_name`: The name of the Google Cloud Run service
* `image`: The container image that will be used for the service
* `delete`: If the service needs to be deleted
* `allow_unauthenticated`: Allow public access to the service

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
type: docker
name: deploy-service

environment: &buildEnv
  service_account_json:
    from_secret: SERVICE_ACCOUNT_JSON
  GOOGLE_CLOUD_PROJECT:
     from_secret: GOOGLE_CLOUD_PROJECT
  GOOGLE_CLOUD_REGION:
     from_secret: GOOGLE_CLOUD_REGION
  
steps:
- name: publish
  image: quay.io/kameshsampath/drone-gcloud-run
  settings:
    service_account_json: ${SERVICE_ACCOUNT_JSON}
    project: ${GOOGLE_CLOUD_PROJECT}
    region: ${GOOGLE_CLOUD_REGION}
    service_name: my-service
    image: asia.gcr.io/${GOOGLE_CLOUD_PROJECT}/greeter
  environment: *buildEnv
```

__IMPORTANT__: It is highly recommended that the environment variables are passed using secrets e.g. `drone exec --secret-file=.env`

Please check the [examples](./examples/) directory for more examples.

## Building

Build the plugin binary:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t quay.io/kameshsampath/drone-gcloud-run -f docker/Dockerfile .
```

## Testing

Execute the plugin from your current working directory:

```text
docker run --rm -e PLUGIN_SERVICE_ACCOUNT_JSON=foo \
  -e PLUGIN_GOOGLE_CLOUD_PROJECT=bar \
  -e PLUGIN_GOOGLE_CLOUD_REGION=asia-south1 \
  -e PLUGIN_IMAGE=asia.gcr.io/${GOOGLE_CLOUD_PROJECT}/greeter \
  -e PLUGIN_SERVICE_NAME=my-service \
  -w /drone/src \
  -v $(pwd):/drone/src \
  quay.io/kameshsampath/drone-gcloud-run
```
