class Route():
    path = ""
    method = ""
    handle_func_name = ""
    handle_func = None

    def __init__(self, path, method, handle_func):
        self.path = path
        self.method = method
        self.handle_func = handle_func
        self.handle_func_name = handle_func.__name__

    def __dict__(self):
        return {
            "path": self.path,
            "method": self.method,
            "handle_func": self.handle_func_name
        }