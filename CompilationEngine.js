fs = require('fs');
buffer = [];
tokens = [];
output = [];
var token;

op = ["+", "-", "*", "/", "&". "|". "<", ">", "="];
keyword_constant = ["true", "false", "null", "this"];
unaryop = ["-", "~"];

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

function checkToken() {
    var args = Array.prototype.slice.call(arguments);
    if (args.indexOf(token.token) == -1 ) { return false };
    return true;
};

function checkTokenType() {
    var args = Array.prototype.slice.call(arguments);
    if (args.indexOf(token.token_type) == -1 ) { return false };
    return true;
};

function checkNextToken() {
    var args = Array.prototype.slice.call(arguments);
    if (args.indexOf(tokens[0].token)) { return false };
    return true;
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
        advance();
        return true;
    } else {
        console.error("missing end of line semicolon");
        return false;
    }   
};

function checkTypeAndIdentifier() {
    //check type
    if (!checkToken("int", "char", "boolean") && !checkIdentifier()) {
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

function checkOpeningBracket() {
    if (checkToken("(")) {
        writeToken();
        advance();
        return true;
    } else {
        console.error("missing opening bracket");
        return false;
};

function checkClosingBracket() {
    if (checkToken(")")) {
        writeToken();
        advance();
        return true;
    } else {
        console.error("missing closing bracket");
        return false;
};

function checkOpeningBrace() {
    if (checkToken("{")) {
        writeToken();
        advance();
        return true;
    } else {
        console.error("missing opening brace");
        return false;
};

function checkClosingBrace() {
    if (checkToken("}")) {
        writeToken();
        advance();
        return true;
    } else {
        console.error("missing closing brace");
        return false;
};

function checkSubroutineCall() {
    
    if (!checkIdentifier()) { return false };
    
    writeToken();
    advance();

    //could be subroutineName(blargh) | class/var.subroutine.(blargh)

    if (checkToken(".")) {
        writeToken();
        advance();
        
        if (!checkIdentifier()) { return false };
        
        writeToken();
        advance();
    }
    
    if (!checkOpeningBracket()) { return false };
    if (!compileExpressionList()) { return false };
    if (!checkClosingBracket()) { return false };

    return true;
};

function compileTerm() {

    ///int, string or keyword
    if (checkTokenType("integerConstant", "stringConstant", "keywordConstant")) { 
        writeToken();
        advance();
        return true;
    }

    //(expression)
    if (checkOpeningBracket()) {
        if (!compileExpression()) { return false };
        if (!checkClosingBracket()) { return false }; 
        return true;
    }

    //unaryOp term    
    if (unaryop.indexOf(token.token) != -1) {
        writeToken();
        advance();

        if (!compileTerm()) { return false };
        return true;
    }

    //can now only have varName | varName [expression] | subroutineCall
    //All of these terms begin with a varName ==> look ahead one token to differentiate


    if (checkNextToken("[")) {
        writeToken();
        advance();

        if (!compileExpression()) { return false };
        if (!checkToken("]")) { return false };

        writeToken();
        advance();

        return true;
    }

    if (checkNextToken("(")) {
        if (!checkSubroutineCall()) { return false };
        return true;
    }

    if (!checkIdentifier()) { return false } 
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
        while (checkToken("constructor") || checkToken("function") || checkToken("method")) {
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
    if (!checkToken("void", "int", "char", "bool")  || !checkIdentifier()) {
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
    
    //possible varDec
    while (checkToken("var")) {  
        if (!compileVarDec()) { return false };
    }
   
    //compile statements
    if (!compileStatements()) { return false };

    //closing brace
    if (!checkToken("}")) { return false };
    
    writeToken();
    advance();

    return true;
};

function compileParameterList() {
    //no paramters
    if (checkToken(")")) { return true };

    if (!checkTypeAndIdentifier()) { return false };

    while(checkToken(",")) {
        if (!checkTypeAndIdentifier()) { return false };
    }
    
    return true;
};

function compileVarDec() {
    writeToken();
    advance();

    if (!checkTypeAndIdentifier()) { return false };

    while (checkToken(",")) {
        if (!checkIdentifier()) { return false };
        writeToken();
        advance();
    }

    if (!checkSemicolon()) { return false };
    
    return true;
};

function compileStatements() {
    while (checkToken("let", "if", "while", "do", "return")) {

        switch(token.token) {
            case "let":
                if (!compileLet()) { return false };
                break;
            case "if":
                if (!compileIf()) { return false };
                break;
            case "while":
                if (!compileWhile()) { return false };
                break;
            case "do":
                if (!compileDo()) { return false };
                break;
            case "return":
                if (!compileReturn()) { return false };
                break;
            default:
                console.error("invalid statement keyword");
                return false;
        }
    }
    return true;
};

function compileLet() {
    writeToken();
    advance();

    //varName 
    if (!checkIdentifier()) { return false };

    writeToken();
    advance();

    //potential expression
    if (checkToken("[")) {
        writeToken();
        advance();
    
        compileExpression();

        if (!checkToken("]")) { return false };
        
        writeToken();
        advance();
    }

    //equals expression
    if (!checkToken("=")) { return false };

    writeToken();
    advance();

    compileExpression();

    if (!checkSemicolon()) { return false }; 
    
    return true;
};

function compileIf() {
    writeToken();
    advance();

    //varName
    if (!checkIdentifier()) { return false };

    writeToken();
    advance();

    //expression
    if (!checkOpeningBracket()) { return false };

    writeToken();
    advance();

    if (!compileExpression()) { return false };

    if (!checkClosingBracket()) { return false };

    if (!checkOpeningBrace()) { return false };
    
    if (!compileStatements()) { return false };

    if (!checkClosingBrace()) { return false };

    if (checkToken("else")) {
        if (!checkOpeningBrace()) { return false };
        if (!compileStatements()) { return false };
        if (!checkClosingBrace()) { return false };
    }
   
    return true
};

function compileWhile() {
    writeToken();
    advance();

    if (!checkOpeningBracket()) { return false };

    if (!checkExpression()) { return false };
    
    if (!checkClosingBracket()) { return false };

    if (!checkOpeningBrace()) { return false };

    if (!compileStatements()) { return false };

    if (!checkClosingBrace()) { return false };
    
    return true;
};

function compileDo() {
    writeToken();
    advance();

    //subroutine call
    if (!checkIdentifier()) { return false };

    if (!checkOpeningBracket()) { return false };

    
    if (!checkClosingBracket()) { return false };

    if (!checkSemicolon()) { return false };

    return true;
};

function compileReturn() {
    writeToken();
    advance();

    if (!compileExpression()) { return false };

    return true;
};

function compileExpression() {
    //an expression must contain at least one term
    if (!compileTerm) { return false };

    //(op term)*
    while (op.indexOf(token.token) != -1) {
        writeToken();
        advance();

        if (!compileTerm()) { return false };
    }
    return true;
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
