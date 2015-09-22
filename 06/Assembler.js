fs = require('fs');
var buffer = [];
var output = [];
var empty = (1 << 16).toString(2).slice(1, 17);
var path = process.argv[2];

var destTable = {
            'M': 1 << 3,
            'D': 1 << 4,
            'A': 1 << 5
};

var jumpTable = {
            'null': 0,
            'JGT': 1,
            'JEQ': 2,
            'JGE': 3,
            'JLT': 4,
            'JNE': 5,
            'JLE': 6,
            'JMP': 7
};

var compTable = {
            '0':   42 << 6,  // 0000101010000000
            '1':   63 << 6,  // 0000111111000000
            '-1':  58 << 6,  // 0000111010000000
            'D':   12 << 6,  // 0000001100000000
            'A':   48 << 6,  // 0000110000000000
            '!D':  13 << 6,  // 0000001101000000
            '!A':  49 << 6,  // 0000110001000000
            '-D':  15 << 6,  // 0000001111000000
            '-A':  51 << 6,  // 0000110011000000
            'D+1': 31 << 6,  // 0000011111000000
            'A+1': 55 << 6,  // 0000110111000000
            'D-1': 14 << 6,  // 0000001110000000
            'A-1': 50 << 6,  // 0000110010000000
            'D+A':  2 << 6,  // 0000000010000000
            'D-A': 19 << 6,  // 0000010011000000
            'A-D':  7 << 6,  // 0000000111000000
            'D&A':  0 << 6,  // 0000000000000000
            'D|A': 21 << 6,  // 0000010101000000
};

function commandType(command) {
    if (command[0] == "@") {
        return (0).toString(2);
    } else {
        return (7 << 13).toString(2);
    }
};

function pad(binary) {
    while (binary.length < 16) {
        binary = '0' + binary;
    };
    return binary;
};

//read file specified in the command line 
data = fs.readFileSync(process.argv[2], "ascii"); 
//generate an array consisting of the commands 
buffer = data.split(/\n/); 
//split each line by commands and symbols 
for (var i=0; i<buffer.length; i++) {
    if (buffer[i][0] == "@") {
        buffer[i] = buffer[i].split(/(@)/);
    }
    else if (buffer[i][1] == "=") {
        buffer[i] = buffer[i].split(/([MDA])+(=)/);
    } else if (buffer[i][2] == "=") {
        buffer[i] = buffer[i].split(/([MDA])+([MDA])+(=)/);
    } else if (buffer[i][3] == "=") {
        buffer[i] = buffer[i].split(/([MDA])+([MDA])+([MDA])+(=)/);
    } else if (buffer[i][1] == ";") {
        buffer[i] = buffer[i].split(/([MDA01])+(;)/);
    } else {
        buffer[i] = buffer[i].split(/(\/)/);
    };
    for (var j=0; j<buffer[i].length; j++) {
        if (buffer[i][j] == "") {
            buffer[i].splice(j, 1);
        };
    };
};
buffer.splice((buffer.length-1), 1);

for (var i=0; i<buffer.length; i++) {
    //is the line a comment?
    if (buffer[i][0] && buffer[i][1] == "/") {
        continue;
    };
    if (buffer[i].length == 0) {
        continue;
    };
    //a or c instruction, generate binary
    binary = commandType(buffer[i]);
    //if a instruction
    if (binary[0] == 0) {
        //1000000000000000 | number
        binary = (parseInt(binary, 2) | buffer[i][1]).toString(2);
        binary = pad(binary);
        output.push(binary);
        continue;
    };
    if (buffer[i][1] == "=") {
        //destination
        binary = (parseInt(binary, 2) | destTable[buffer[i][0]]).toString(2);
        //a-bit
        if(buffer[i][2].indexOf('M') != -1) {
            binary = (1 << 12 | parseInt(binary, 2)).toString(2);
            //must contain M. Messy hack to grab the binary from the comp table
            compute = (buffer[i][2]).replace(/M/g, "A");
            binary = (parseInt(binary, 2) | compTable[compute]).toString(2);
        } else {
            binary = (parseInt(binary, 2) | compTable[buffer[i][2]]).toString(2);
        };
    } else {
        // ';' 
        binary = (parseInt(binary, 2) | compTable[buffer[i][0]]).toString(2);
        binary = (parseInt(binary, 2) | jumpTable[buffer[i][2]]).toString(2);
    };
    output.push(binary);
};
output = output.join("\n");
var fd =fs.openSync(path + ".hack", "w");
fs.write(fd, output);
