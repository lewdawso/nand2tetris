re('fs');
var buffer = [];
var binary = [];
var empty = (1 << 16).toString(2).slice(1, 17);

function commandType(command) {
    if (command[0] == "@") {
        return (1 << 15).toString(2);
    } else {
        return (7 << 13).toString(2);
    }
};

//read file specified in the command line data = fs.readFileSync("test", "ascii"); //generate an array consisting of the commands buffer = data.split(/\n/); //split each line by commands and symbols for (var i=0; i<buffer.length; i++) {
    if (buffer[i][0] == "@") {
        buffer[i] = buffer[i].split(/(@)/);
    }
    else {
        buffer[i] = buffer[i].split(/([MDA])+=/);
    };
    for (var j=0; j<buffer[i].length; j++) {
        if (buffer[i][j] == "") {
            buffer[i].splice(j, 1);
        };
    };
};
//get rid of empty array at end of buffer //buffer.splice((buffer.length-1), 1);

for (var i=0; i<buffer.length; i++) {
    command = commandType(buffer[i]);
    if (command[0] == 
    binary.push(command);
};

console.log(buffer);
console.log(binary);

