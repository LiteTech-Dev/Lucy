import { ModrinthApiHandler } from "../networkUtils/modrinthApiHandler.js";
export async function runTest() {
    ModrinthApiHandler.getInstance()
        .getProject("betterfps")
        .then(function (result) {
            console.log(result);
        })
        .catch(function (err) {
            console.error(err);
        });
}

runTest();
