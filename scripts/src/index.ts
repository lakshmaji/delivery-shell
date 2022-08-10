import Ajv from "ajv"
import * as fs from 'fs';
import * as path from 'path';
import schema from "./schema";


const main = () => {
  const ajv = new Ajv()
  const validate = ajv.compile(schema)
  
  const encodedString = fs.readFileSync(path.join(__dirname, "../../offers.json" ), 'utf8');

  if (validate(JSON.parse(encodedString))) {
    console.log("Valid configuration")
  } else {
    console.error("Invalid configuration")
    console.log(validate.errors)
  }
  
}

main()