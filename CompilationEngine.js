fs = require('fs');
buffer = [];
tokens = [];
output = [];
var token;
var previousToken;

op = ["+", "-", "*", "/", "&", "|", "<", ">", "="];
keyword_constant = ["true", "false", "null", "this"];
unaryop = ["-", "~"];

function debug() {
    console.log("*** DEBUG ***");
    console.log("previous token: " + previousToken.token);
    console.log("current token: " + token.token);
    console.log("next token: " + tokens[0].token);
};

function write(cmd) {
	output.push(cmd);
};

function writeOpen(cmd) {
    output.push("<" + cmd + ">")
}

function writeClose(cmd) {
    output.push("</" + cmd + ">")
}

function writeToken() {
    write(["<" + token.token_type + ">" + " " + token.token + " " + "</" + token.token_type + ">" ]);    
};

function genOutFile() {
	output = output.join("\n");
	path = process.argv[2].split("/")
    name = process.argv[3]
	fd = fs.openSync(path[0] + "/" + path[1] + "/" + name + "_.xml", 'w');
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
    previousToken = token;
    token = tokens.shift();
};

function generateArray(buffer) {
    for (var i=0; i<buffer.length; i++) {
        temp = buffer[i].split(/"/).filter(function(e){return e!="" && e!=null && e!=" "});
        tokens.push({"token_type":temp[0], "token":temp[1]});
    }
};

function raiseError(message) {
    console.error("error: " + message);
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
    if (args.indexOf(tokens[0].token) == -1) { return false };
    return true;
};

function checkIdentifier() {
    if ((token.token).match(/^[a-zA-Z_][\w]*$/)) { 
        return true;
    } else {
        raiseError("invalid identifier");
        return false;
    }
};

function checkSemicolon() {
    if (checkToken(";")) {
        writeToken();
        advance();
        return true;
    } else {
        raiseError("missing end of line semicolon");
        return false;
    }   
};

function checkTypeAndIdentifier() {
    //check type
    if (!checkToken("int", "char", "boolean") && !checkIdentifier()) {
        raiseError("missing type specifier in var dec");
        return false;
    }
    
    writeToken();
    advance();

    //check variable name
    if (!checkIdentifier) { raiseError("var dec identifier is invalid") ; return false };
    
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
        raiseError("missing opening bracket");
        return false;
    }
};

function checkClosingBracket() {
    if (checkToken(")")) {
        writeToken();
        advance();
        return true;
    } else {
        raiseError("missing closing bracket");
        return false;
    }
};

function checkOpeningBrace() {
    if (checkToken("{")) {
        writeToken();
        advance();
        return true;
    } else {
        raiseError("missing opening brace");
        return false;
    }
};

function checkClosingBrace() {
    if (checkToken("}")) {
        writeToken();
        advance();
        return true;
    } else {
        raiseError("missing closing brace");
        return false;
    }
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
    
    writeOpen("term")
    ///int, string or keyword
    if (checkTokenType("integerConstant", "stringConstant", "keywordConstant")) { 
        writeToken();
        advance();
        return true;
    }

    //(expression)
    if (checkToken("(")) {
        if (!compileExpression()) { raiseError("compileExpression") ; return false };
        if (!checkClosingBracket()) { return false }
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
        if (!checkIdentifier()) { return false };
        writeToken();
        advance();

        //write "["
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
    writeToken();
    advance();
    writeClose("term")
    return true;
};

function compileClass() {
    advance();
    if (checkToken("class")) {
        writeOpen("class")
        writeToken();
        advance();
        if (!checkIdentifier()) { raiseError("invalid class identifier"); return false };   
        writeToken();
        advance();
        if (!checkToken("{")) {
            raiseError("no opening brace for class");
            return false;
        }
        writeToken();
        advance();
        
        //classVarDec
        while (checkToken("static") || checkToken("field")) {
            if (!compileClassVarDec()) { return false };
        }

        //subroutineDec
        while (checkToken("constructor") || checkToken("function") || checkToken("method")) {
            if (!compileSubroutine()) { raiseError("compileSubroutine"); return false };
        }

        if (!checkClosingBrace()) { return false }

        if (tokens.length != 0) {
            raiseError("tokens present after class closed");
            return false;
        }
        writeClose("class")
        return true;

    } else {
        raiseError("missing class keyword to open");
        return false;
    }
};

function compileClassVarDec() {
    writeOpen("classVarDec")
    writeToken();
    advance();
    //check type
    if (!checkToken("int", "char", "bool") && !checkIdentifier()) {
        raiseError("missing type specifier");
        return false;
    }
    writeToken();
    advance();
    if (!checkIdentifier()) { raiseError("checkIdentifier") ; return false };
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
    
    writeClose("classVarDec")
    return true;
}

function compileSubroutine() {
    writeOpen("subroutineDec")
    writeToken();
    advance();

    //check return type
    if (!checkToken("void", "int", "char", "bool")  && !checkIdentifier()) {
        raiseError("missing subroutine return type");
        return false;
    }
    
    writeToken();
    advance();

    //check subroutine name
    if (!checkIdentifier()) { raiseError("invalid subroutine identifier") ; return false };
    writeToken();
    advance();

    //opening bracket
    if (!checkOpeningBracket()) { return false };

    //parameter list
    if (!compileParameterList()) { return false };
    
    //closing bracket
    if (!checkClosingBracket()) { return false };

    writeOpen("subroutineBody");

    //opening brace
    if (!checkOpeningBrace()) { return false };
    
    //possible varDec
    while (checkToken("var")) {  
        if (!compileVarDec()) { 
            raiseError("error: compileVarDec");
            return false;
        }
    }

    //compile statements
    if (!compileStatements()) { 
        raiseError("compileStatements");
        return false;
    }

    //closing brace
    if (!checkClosingBrace()) { return false };

    writeClose("subroutineBody");

    writeClose("subroutineDec");
    return true;
};

function compileParameterList() {
    //no paramters
    writeOpen("parameterList")
    if (checkToken(")")) { writeClose("parameterList") ; return true };

    if (!checkTypeAndIdentifier()) { raiseError("checkTypeAndIdentifier") ; return false };

    while(checkToken(",")) {
        if (!checkTypeAndIdentifier()) { return false };
    }
    
    writeClose("parameterList")
    return true;
};

function compileVarDec() {
    writeOpen("varDec")
    writeToken();
    advance();

    if (!checkTypeAndIdentifier()) { raiseError("checkTypeAndIdentifier") ; return false };

    while (checkToken(",")) {
        if (!checkIdentifier()) { return false };
        writeToken();
        advance();
    }

    if (!checkSemicolon()) { return false };
    
    writeClose("varDec")
    return true;
};

function compileStatements() {
    writeOpen("statements")
    while (checkToken("let", "if", "while", "do", "return")) {

        switch(token.token) {
            case "let":
                if (!compileLet()) { raiseError("compileLet") ; return false };
                break;
            case "if":
                if (!compileIf()) { raiseError("compileIf") ; return false };
                break;
            case "while":
                if (!compileWhile()) { raiseError("compileWhile") ; return false };
                break;
            case "do":
                if (!compileDo()) { raiseError("compileDo") ; return false };
                break;
            case "return":
                if (!compileReturn()) { raiseError("compileReturn") ; return false };
                break;
            default:
                raiseError("invalid statement keyword");
                return false;
        }
    }
    writeClose("statements")
    return true;
};

function compileLet() {
    writeOpen("letStatement")
    writeToken();
    advance();

    //varName 
    if (!checkIdentifier()) { raiseError("invalid compileLet identifier") ; return false };

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
    if (!checkToken("=")) { raiseError("missing =") ; return false };

    writeToken();
    advance();

    if (!compileExpression()) { raiseError("unable to compileExpression inside compileLet") ; return false }

    if (!checkSemicolon()) { return false }; 
    
    writeClose("letStatement");
    return true;
};

function compileIf() {
    writeOpen("ifStatment");
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
    writeClose("ifStatement");
    return true
};

function compileWhile() {
    writeOpen("whileStatement");
    writeToken();
    advance();

    if (!checkOpeningBracket()) { return false };

    if (!compileExpression()) { return false };
    
    if (!checkClosingBracket()) { return false };

    if (!checkOpeningBrace()) { return false };

    if (!compileStatements()) { return false };

    if (!checkClosingBrace()) { return false };
    
    writeClose("whileStatement");
    return true;
};

function compileDo() {
    writeOpen("doStatement");
    writeToken();
    advance();

    if (!checkSubroutineCall()) { return false };
    if (!checkSemicolon()) { return false };
    writeClose("doStatement");
    return true;
};

function compileReturn() {
    writeOpen("returnStatement");
    writeToken();
    advance();
    //expression is optional
    if (checkToken(";")) {
        writeToken();
        advance();
        writeClose("returnStatement")
        return true;
    }

    if (!compileExpression()) { raiseError("compileExpression") ; return false };
    if (!checkSemicolon()) { return false };
    writeClose("returnStatement");
    return true;
};

function compileExpression() {
    writeOpen("expression");
    //an expression must contain at least one term
    if (!compileTerm()) { raiseError("compileTerm") ; return false };

    //(op term)*
    while (op.indexOf(token.token) != -1) {
        writeToken();
        advance();
        if (!compileTerm()) { return false };
    }
    writeClose("expression");
    return true;
};

function compileExpressionList() {
    writeOpen("expressionList");
    //first off, check if we have an empty list
    if (checkToken(")")) { writeClose("expressionList") ; return true };
    
    //now we know we have to compile at least one expression
    if (!compileExpression()) { return false };

    while (!checkToken(")")) { 
        if (!checkToken(",")) { return false };
        writeToken();
        advance();
        if (!compileExpression()) { return false };
    }
    writeClose("expressionList");
    return true;
};

function main() {
    file = process.argv[2] 
    data = fs.readFileSync(file, 'ascii');
    buffer = data.split(/[\n]/);
    generateArray(buffer);
    
    //first routine to be called must be compileClass
    if (!compileClass()) { raiseError("unable to compile class") ; debug() ; return }
	genOutFile();
};

main();
