fs = require('fs');
buffer = [];
tokens = [];
output = [];
var token;

function write(cmd) {
	output.push(cmd)
};

function writeToken(token_type, token) {
    write(["<" + token_type + ">" + " " + token + " " + "</" + token_type + ">" ]);    
};

function genOutFile() {
	output = output.join("\n");
	name = process.argv[2].split(".");
	fd = fs.openSync(name[0] + "_.xml", 'w');
	fs.write(fd, output);
};

function hasMoreTokens() {
    if(tokens.length) {
        return true;
    } else {
        return false;
    };
};

function advance() {
    return tokens.shift();
};

function generateArray(buffer) {
    for (var i=0; i<buffer.length; i++) {
        temp = buffer[i].split(/"/).filter(function(e){return e!="" && e!=null && e!=" "});
        tokens.push({"token_type":temp[0], "token":temp[1]});
    }
};

function checkToken(token, expect) {
    if (token == expect) {
        return true;
    }
    return false;
};

function checkIdentifier(identifier) {
    if (identifier.match(/^[a-zA-Z_][\w]*$/)) { 
        return true;
    }
    return false;
};

function compileClass() {
    token = advance();
    if (checkToken(token.token, "class")) {
        write(["<class>"]);
        writeToken(token.token_type, token.token);
        token = advance();
        if (checkIdentifier(token.token)) {
            writeToken(token.token_type, token.token);
        } else {
            console.error("invalid class identifier");
            return;
        }
        token = advance();
        if (!checkToken(token.token, "{")) {
            console.error("no opening brace for class");
            return;
        }
        writeToken(token.token_type, token.token);
        token = advance();
        
        //classVarDec
        if (checkToken(token.token, "static") || checkToken(token.token, "field")) {
            compileClassVarDec();
        }

        //subroutineDec
        if (checkToken(token.token, "constructor"), checkToken(token.token, "function"), checkToken(token.token, "method")) {
            compileSubroutine();
        }
    }  
};

function compileClassVarDec() {
    //stub
    return;
};

function compileSubroutine() {
    //stub
    return;
};

function main() {
    file = process.argv[2] 
    data = fs.readFileSync(file, 'ascii');
    buffer = data.split(/[\n]/);
    generateArray(buffer);
    
    //first routine to be called must be compileClass
    compileClass();
	//genOutFile();
};

main();
