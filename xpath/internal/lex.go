package internal

type LexKind uint

const (
	LexComma      LexKind = iota //= ','
	LexSlash                     // '/',
	LexAt                        //'@',
	LexDot                       // '.',
	LexLParens                   // '(',
	LexRParens                   // ')',
	LexLBracket                  // '[',
	LexRBracket                  // ']',
	LexStar                      // '*',
	LexPlus                      // '+',
	LexMinus                     // '-',
	LexEq                        // '=',
	LexLt                        // '<',
	LexGt                        // '>',
	LexBang                      // '!',
	LexDollar                    // '$',
	LexApos                      // '\'',
	LexQuote                     // '"',
	LexUnion                     // '|',
	LexNe                        // 'N',   // !=
	LexLe                        // 'L',   // <=
	LexGe                        // 'G',   // >=
	LexAnd                       // 'A',   // &&
	LexOr                        // 'O',   // ||
	LexDotDot                    // 'D',   // ..
	LexSlashSlash                // 'S',   // //
	LexName                      // 'n',   // XML _Name
	LexString                    // 's',   // Quoted string constant
	LexNumber                    // 'd',   // _Number constant
	LexAxe                       // 'a',   // Axe (like child::)
	LexEof                       // 'E',
)
