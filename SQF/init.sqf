_result = "a3Logs" callExtension ["atm", ["playerId", "123321", "action", "deposit", "amount", "999999", "bank", "1000000", "hand", "999999"]];
diag_log format ["a3Logs: Result: %1", _result];
_result = "a3Logs" callExtension ["atm", ["playerId", "123321", "action", "withdraw", "amount", "999999", "bank", "1000000", "hand", "999999"]];
diag_log format ["a3Logs: ATM withdraw: %1", _result];
_result = "a3Logs" callExtension ["sell", ["playerId", "123321", "resource", "t_peach", "amount", "10", "sum", "100", "hand", "1000"]];
diag_log format ["a3Logs: Sell: %1", _result];
