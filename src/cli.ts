import { Command } from "commander";
import { FileSystem } from "./backend.js";

const program = new Command();

program.name("mcpm").description("");

program
    .command("init")
    .description("")
    .action(() => {
        const fileSystem = new FileSystem();

        console.log(fileSystem.extractFromJar("server.jar", "version.json"));
    });

program.command("list");

program.command("install");

program.command("remove");

program.command("update");

program.command("disable");

program.command("enable");

program.command("export");

program.command("import");

program.command("host");

export { program };
