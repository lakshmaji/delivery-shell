import {JSONSchemaType} from "ajv"
import { Offers } from "./models";

const schema: JSONSchemaType<Offers> = {
    type: "array",
    items: {
        type: "object",
        properties: {
        code: {
            type: "string"
        },
        discount: {
            type: "number"
        },
        conditions: {
            type: "array",
            items: {
            type: "object",
            properties: {
                fact: {
                type: "string",
                enum: ["distance", "weight"]
                },
                operator: {
                type: "string",
                enum: ["lessThan", "greaterThanOrEqual", "lessThanOrEqual"]
                },
                value: {
                type: "number"
                }
            },
            required: ["fact", "operator", "value"],
            additionalProperties: false,
            },
            minItems: 1,
            maxItems: 30,
        },
        },
        required: ["code", "discount", "conditions"],
        additionalProperties: false,
    },
    minItems: 1,
    maxItems: 50,
}

export default schema