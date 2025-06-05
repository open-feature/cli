import { OpenFeature } from "@openfeature/server-sdk";
import generatedClient from "./generated/index.js";

async function main(){
    try {
        console.log('Testing generated OpenFeature Node.js ...')

        await OpenFeature.setProvider(generatedClient)

        const client = OpenFeature.getClient()
    } catch (error) {
        
    }
}