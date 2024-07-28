import { ModrinthApiHandler } from "../networkUtils/modrinthApiHandler.js";
export async function runTest() {
    ModrinthApiHandler.getInstance()
        .getVersionByHash("2bece942d05315e512b468301523136c1d79d8b7")
        .then(function (result) {
            console.log(result);
        })
        .catch(function (err) {
            console.error(err);
        });
}

runTest();
