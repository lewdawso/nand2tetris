fs = require('fs');
trimmed = [];
parsed = [];
tokens = [];
output = [];

keywords = 
[
    "class",
    "constructor",
    "function",
    "method",
    "field",
    "static",
    "var",
    "int",
    "char",
    "boolean",
    "void",
    "true",
    "false",
    "null",
    "this",
    "let",
    "do",
    "if",
    "else",
    "while",
    "return"
];

symbols = ["{", "}", "(", ")", "[", "]", ".", ",", ";", "+", "-", "*", "/", "&", "|", "<", ">", "=", "~"];

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
        if (command[i] == "/") {
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
        trimmed[j] = trimmed[j].split(/([{()};.])|[, ]/);
    };
    for (var i=0; i<trimmed.length; i++) {
        command = trimmed[i]; 
        for (var j=0; j<command.length; j++) {
            if (command[j] != "") {
                parsed.push(command[j]);
            };    
        };
    };   
    //remove whitespace
    for (var i=0; i<parsed.length; i++) {
        if (!parsed[i] == "") {
            tokens.push(parsed[i]);
        };
    };
};

function write(cmd) {
	output.push(cmd)
}

function genOutFile() {
	output = output.join("\n");
	name = process.argv[2].split(".");
	fd = fs.openSync(name[0] + "_T.xml", 'w');
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
}

function tokenType(token) {
    if (keywords.indexOf(token) != -1) {
        return "keyword";
    } else if (symbols.indexOf(token) != -1) {
		return "symbol";
    } else if (token <=0 && token <=32767) {
		return "int_const";
    } else if (token[0] ==  '"' && token[token.length-1] == '"' ) {
		return "string_const";
    } else {
		return "identifier";
    };
};

function main() {
    file = process.argv[2] 
    data = fs.readFileSync(file, 'ascii');
    buffer = data.split(/[\n\r\t]/);
    parse();
    var token;
    var token_type;
	var token_value;
	write(["<tokens>"]);
	while(hasMoreTokens()) {
        token = advance(tokens);
    	token_type = tokenType(token);	
		if (token_type == "string_const") {
			token.replace('"', '');
		};
		if (token_type == "symbol") {
			switch(token) {
				case "<":
					token = "&lt";
					break;
				case ">":
					token = "&gt";
					break;
				case '"':
					token = "&quot";
					break;
				case "&":
					token = "&amp";
					break;
			};
		};
		write(["  <" + token_type + ">" + " " + token + " " + "</" + token_type + ">" ]);
	};
	write(["</tokens>"]);
	genOutFile();
};

main();
