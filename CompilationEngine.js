fs = require('fs');
buffer = [];
tokens = [];
output = [];
var token;

function write(cmd) {
	output.push(cmd)
};

function writeToken() {
    write(["<" + token.token_type + ">" + " " + token.token + " " + "</" + token.token_type + ">" ]);    
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
    token = tokens.shift();
};

function generateArray(buffer) {
    for (var i=0; i<buffer.length; i++) {
        temp = buffer[i].split(/"/).filter(function(e){return e!="" && e!=null && e!=" "});
        tokens.push({"token_type":temp[0], "token":temp[1]});
    }
};

function checkToken(expect) {
    if (token.token == expect) {
        return true;
    }
    return false;
};

function checkIdentifier() {
    if ((token.token).match(/^[a-zA-Z_][\w]*$/)) { 
        return true;
    } else {
        console.error("invalid identifier"); 
        return false;
    }
};

function checkSemicolon() {
    if (checkToken(";")) {
        writeToken();
        return true;
    } else {
        console.error("missing end of line semicolon");
        return false;
    }   
}

function checkTypeAndIdentifier() {
    //check type
    if (!checkToken("int") || !checkToken("char") || !checkToken("bool") || !checkIdentifier()) {
        console.error("missing type specifier");
        return false;
    }
    
    writeToken();
    advance();

    //check variable name
    if (!checkIdentifier) { return false };
    
    writeToken();
    advance();
    return true;
};

function compileClass() {
    advance();
    if (checkToken("class")) {
        write(["<class>"]);
        writeToken();
        advance();
        if (!checkIdentifier()) { return false };   
        writeToken();
        advance();
        if (!checkToken("{")) {
            console.error("no opening brace for class");
            return;
        }
        writeToken();
        advance();
        
        //classVarDec
        while (checkToken("static") || checkToken("field")) {
            compileClassVarDec();
        }

        //subroutineDec
        while (checkToken("constructor"), checkToken("function"), checkToken("method")) {
            compileSubroutine();
        }

        advance();
        if (!checkToken("}")) {
            console.error("no closing brace for class");
            return;
        }

        advance();
        if (tokens.length != 0) {
            console.error("tokens present after class closed");
            return;
        }

    } else {
        console.error("missing class keyword to open");
        return;
    }
};

function compileClassVarDec() {
    writeToken();
    advance();
    
    //check type
    if (!checkToken("int") || !checkToken("char") || !checkToken("bool") || !checkIdentifier()) {
        console.error("missing type specifier");
        return false;
    }
    writeToken();
    advance();
    if (!checkIdentifier()) { return false };
    writeToken();
    advance();
    
    //deal with (',', varName)*
    while (token.token == ",") {
        writeToken();
        advance();
        if (!checkIdentifier()) { return false };
        writeToken();
        advance(); 
    }

    if (!checkSemicolon()) { return false };
    writeToken()
    advance();
    
    return true;
}

function compileSubroutine() {
    writeToken();
    advance();

    //check return type
    if (!checkToken("void") || !checkToken("int") || !checkToken("char") || !checkToken("bool") || !checkIdentifier()) {
        console.error("missing subroutine return type");
        return false;
    }
    
    writeToken();
    advance();

    //check subroutine name
    if (!checkIdentifier()) { return false };
    writeToken();
    advance();

    //opening bracket
    if (!checkToken("(")) { return false };
    writeToken();
    advance();

    //parameter list
    if (!compileParameterList()) { return false };
    
    //closing bracket
    if (!checkToken(")")) { return false };
    writeToken();
    advance();

    //subroutine body

    //opening brace
    if (!checkToken("{")) { return false};
    writeToken();
    advance();

    //
    if (checkToken("var")) {  
        if (!compileVarDec()) {return false };
    }
    
    if (!compileStatements()) { return false };

    //closing brace
    if (!checkToken("}")) { return false };
    writeToken();
    advance();

    return true;
};

function compileParameterList() {
    var success;

    //no paramters
    if (checkToken(")")) { return true };
    
    if (!checkTypeAndIdentifier()) { return false };

    while(checkToken(",")) {
        if (!checkTypeAndIdentifier()) {success = false ; break};
    }
    
    return true;
};

function compileVarDec() {
    //stub
    return;
};

function compileStatements() {
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
