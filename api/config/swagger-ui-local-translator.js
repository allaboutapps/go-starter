// https://github.com/swagger-api/swagger-ui/blob/master/docker/configurator/translator.js
// Based on commit in 02758b8125dbf38763cfd5d4f91c7c803e9bd0ad
// Nov 08, 2018
// see "AAA change" for applied changes

// Converts an object of environment variables into a Swagger UI config object
const configSchema = require("./variables")

const defaultBaseConfig = {
    url: {
        value: "https://petstore.swagger.io/v2/swagger.json",
        schema: {
            type: "string",
            base: true
        }
    },
    dom_id: {
        value: "#swagger-ui",
        schema: {
            type: "string",
            base: true
        }
    },
    deepLinking: {
        value: "true",
        schema: {
            type: "boolean",
            base: true
        }
    },
    presets: {
        value: `[\n  SwaggerUIBundle.presets.apis,\n  SwaggerUIStandalonePreset\n]`,
        schema: {
            type: "array",
            base: true
        }
    },
    plugins: {
        value: `[\n  SwaggerUIBundle.plugins.DownloadUrl\n]`,
        schema: {
            type: "array",
            base: true
        }
    },
    layout: {
        value: "StandaloneLayout",
        schema: {
            type: "string",
            base: true
        }
    },
    // AAA change START: interspect requests to 8081 and instead pass them to 8080 where our server runs
    requestInterceptor: {
        value: `function (request) {
  request.url = request.url.split(":8081/").join(":8080/");
  console.log("[requestInterceptor] :8080", request);
  return request;
}`,
        schema: {
            type: "function",
            base: true
        }
    }
    // AAA change END
}

function objectToKeyValueString(env, { injectBaseConfig = false, schema = configSchema, baseConfig = defaultBaseConfig } = {}) {
    let valueStorage = injectBaseConfig ? Object.assign({}, baseConfig) : {}
    const keys = Object.keys(env)

    // Compute an intermediate representation that holds candidate values and schemas.
    //
    // This is useful for deduping between multiple env keys that set the same
    // config variable.

    keys.forEach(key => {
        const varSchema = schema[key]
        const value = env[key]

        if (!varSchema) return

        if (varSchema.onFound) {
            varSchema.onFound()
        }

        const storageContents = valueStorage[varSchema.name]

        if (storageContents) {
            if (varSchema.legacy === true && !storageContents.schema.base) {
                // If we're looking at a legacy var, it should lose out to any already-set value
                // except for base values
                return
            }
            delete valueStorage[varSchema.name]
        }

        valueStorage[varSchema.name] = {
            value,
            schema: varSchema
        }
    })

    // Compute a key:value string based on valueStorage's contents.

    let result = ""

    Object.keys(valueStorage).forEach(key => {
        const value = valueStorage[key]

        const escapedName = /[^a-zA-Z0-9]/.test(key) ? `"${key}"` : key

        if (value.schema.type === "string") {
            result += `${escapedName}: "${value.value}",\n`
        } else {
            result += `${escapedName}: ${value.value === "" ? `undefined` : value.value},\n`
        }
    })

    return result.trim()
}

module.exports = objectToKeyValueString