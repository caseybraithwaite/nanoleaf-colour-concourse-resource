# Nanoleaf Colour Concourse Resource

Resource for manipulating Nanoleaf panels from Concourse. 

## Source Configuration

*ip_address* - **required** - the IP address of the nanoleaf device

*api_token* - **required** - the API token associated with the given device

## Params

*power* - **optional** - whether to power the device on. options: `true`, `false`. default `true`

*hue* - **optional** - the hue to set. options: a value between 0 - 360. default `none`

*brightness* - **optional** - the brightness to set. options: a value between 0 - 100. default `none`

## Example Pipelines

See examples/pipeline.yaml.
