number_list = {1, 2, 3}

object_list = {}
table.insert(object_list, {1})
table.insert(object_list, {2})
table.insert(object_list, {3})

dict_list = {}
table.insert(dict_list, {name = "a", n = 1})
table.insert(dict_list, {name = "b", n = 2})

config = {db = {host = "localhost", port = 5432}, name = "app"}

matrix = {{1, 2}, {3, 4}}

cyclic = {}
table.insert(cyclic, cyclic)
