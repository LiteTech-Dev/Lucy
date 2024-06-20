import { Command } from "commander";
import { initialization } from "./initialization.js";

const program = new Command();

program.name("mcpm").description("");

program
    .command("init")
    .description("")
    .action(() => {
        initialization();
    });

program.command("list");

program.command("install");

program.command("remove");

program.command("update");

program.command("disable");

program.command("enable");

program.command("export");

program.command("import");

export { program };
