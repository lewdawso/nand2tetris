fs=require('fs');
var buffer = [];
var output = [];
var asm = [];
var path = process.argv[2];

var commands = {"add":"C_ARITHMETIC",
                "sub":"C_ARITHMETIC",
                "neg":"C_ARITHMETIC",
                "eq":"C_ARITHMETIC",
                "gt":"C_ARITHMETIC",
                "lt":"C_ARITHMETIC",
                "and":"C_ARITHMETIC",
                "or":"C_ARITHMETIC",
                "not":"C_ARITHMETIC",
                "push":"C_PUSH"
                };

function isComment(command) {
    if ((command[0] == "/") && (command[1] == "/")) {
        return true;
    };
};

function isWhiteSpace(command) {
    if (command.length == 0) {
        return true;
    };
};

function commandType(cmd) {
    return commands[cmd]; 
};

function arg1(cmd) {
    if (commandType(cmd[0]) == "C_ARITHMETIC") {
        return cmd[0];
    } else {
        return cmd[1];
    };
};

function arg2(cmd) {
    if (commandType(cmd[0]) == ("C_PUSH" || "C_POP" || "C_FUNCTION" || "C_CALL")) {
        return cmd[2];
    };
};

function parse() {
    //strip comments and trim whitespace
    for (var i=0; i<buffer.length; i++) {
        command = buffer[i]
        if (!isComment(command) && !isWhiteSpace(command)) {
            output.push(buffer[i]);
        };
    };
    //split commands into individual keywords
    for (var j=0; j<output.length; j++) {
        output[j] = output[j].split(" ");
    };
};

function write(cmd) {
    asm.push(cmd);
};

function genRandLabel() {
    var text = "";
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWZYXabcdefghijklmnopqrstuvwxyz";
    for (var i=1; i<6; i++) {
        text +=possible.charAt(Math.floor(Math.random() * possible.length));
    };
    return text;
};

function writeArithmetic(cmd) {
    label1 = genRandLabel();
    label2 = genRandLabel();
    switch(arg1(cmd)) {
        case "add":
            getTop2Stack();
            write(["M=M+D"]);
            incrementRegister("SP");
            break;
        case "sub":
            getTop2Stack();
            write(["M=M-D"]);
            incrementRegister("SP");
            break;
        case "neg":
            decrementRegister("SP");
            setAtoSP();
            write(["M=-M"]);
            incrementRegister("SP");
            break;
        case "eq":
            getTop2Stack();
            write(["D=D-M"]);
            write(["@"+label1]);
            write(["D;JNE"]);
            write(["@1"]);
            write(["D=A"]);
            setAtoSP();
            write(["M=D"]);
            write(["@"+label2]);
            write(["0;JMP"]);
            write(["("+label1+")"]);
            write(["@0"]);
            write(["D=A"]);
            setAtoSP();
            write(["M=D"]);
            write(["("+label2+")"]);
            incrementRegister("SP");
            break;
        case "gt":
            getTop2Stack();
            write(["D=D-M"]);
            write(["@"+label1]);
            write(["D;JGT"]);
            write(["@1"]);
            write(["D=A"]);
            setAtoSP();
            write(["M=D"]);
            write(["@"+label2]);
            write(["0;JMP"]);
            write(["("+label1+")"]);
            write(["@0"]);
            write(["D=A"]);
            setAtoSP();
            write(["M=D"]);
            write(["("+label2+")"]);
            incrementRegister("SP");
            break;
        case "lt":
            getTop2Stack();
            write(["D=D-M"]);
            write(["@"+label1]);
            write(["D;JLT"]);
            write(["@1"]);
            write(["D=A"]);
            setAtoSP();
            write(["M=D"]);
            write(["@"+label2]);
            write(["0;JMP"]);
            write(["("+label1+")"]);
            write(["@0"]);
            write(["D=A"]);
            setAtoSP();
            write(["M=D"]);
            write(["("+label2+")"]);
            incrementRegister("SP");
            break;
        case "and":
            getTop2Stack();
            write(["M=M&D"]);
            incrementRegister("SP");
            break;
        case "or":
            getTop2Stack();
            write(["M=M|D"]);
            incrementRegister("SP");
            break;
        case "not":
            decrementRegister("SP");
            setAtoSP();
            write(["M=!M"]);
            incrementRegister("SP");
            break;
    };
};

function writePushPop(command, cmdType) {
   switch (cmdType) {
        case "C_PUSH":
            write(["@" + arg2(command)]);
            write(["D=A"]);   
            setAtoSP();
            write(["M=D"]);
            incrementRegister("SP");
            break;
    };
};

function getTop2Stack() {
    decrementRegister("SP");
    setAtoSP();
    write(["D=M"]);
    decrementRegister("SP");
    setAtoSP();
};

function decrementRegister(register) {
    write(["@"+register]);
    write(["M=M-1"]);
};

function setAtoSP() {
    write(["@SP"]);
    write(["A=M"]);
};

function incrementRegister(register) {
    write(["@"+register]);
    write(["M=M+1"]);
};

function genOutFile() {
    asm = asm.join("\n");
    path = path.split(".");
    fd = fs.openSync(path[0] + ".asm", "w");
    fs.write(fd, asm);
};

function main() {
    data=fs.readFileSync(path, "ascii");
    buffer = data.split(/\n/);
    parse();
    for (var k=0; k<output.length; k++) {
        command = output[k];
        var cmdType = commandType(command[0])
        if (cmdType == "C_ARITHMETIC") {
            writeArithmetic(command);
        } else if (cmdType == ("C_PUSH" || "C_POP")) {
            writePushPop(command, cmdType);
        };
    };
    genOutFile()
};

main();
