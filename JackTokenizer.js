fs = require('fs');
trimmed = [];
accepted = [];

function isComment(command) {
    if ((command[0] == "/" && command[1] == "/") || (command[0] == "/" && command[1] == "*") || command[1] == "*" || (command[1] == "*" && command[2] == "/")) {
        return true;
    };
};

function isWhiteSpace(command) {
    if (command.length == 0) {
        return true;
    };
};

function isSpaceComment(command) {
    for (var i=0; i<command.length; i++) {
        if (command[i] == "/" && command[i+1] == "/") {
            return true;
        };
    };
};

//gets things into a workable state for tokenizing
function parse() {
    //strip comments and trim whitespace
    for (var i=0; i<buffer.length; i++) {
        command = buffer[i]
        if (!isComment(command) && !isWhiteSpace(command) && !isSpaceComment(command)) {
            trimmed.push(buffer[i]);
        };
    };
    //split commands into individual keywords
    for (var j=0; j<trimmed.length; j++) {
        trimmed[j] = trimmed[j].split(/([{()}; ])/);
    };
    for (var i=0; i<trimmed.length; i++) {
        command = trimmed[i]; 
        for (var j=0; j<command.length; j++) {
            if (command[j] != "") {
                accepted.push(command[j]);
            };    
        };
    };   
    //remove whitespace
    for (var i=0; i<accepted.length; i++) {
        if (accepted[i] == " " || accepted[i] == "") {
            //accepted.splice(i, 1);
            delete accepted[i];
        };
    };
};

function main() {
    file = process.argv[2] 
    data = fs.readFileSync(file, 'ascii');
    buffer = data.split(/[\n\r\t]/);
    parse();
    console.log(accepted);
}

main();






