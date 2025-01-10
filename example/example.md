# Hello World

This is test 1

```bash {"name": "test1"}
echo "Hello test1"
```

This is test 2

```bash {"name": "test2", "description": "This is a test 2"}
echo "Hello from test2"
read -p "Enter your name: " name
echo "Hello $name"

env | grep -i $name || true
```

This is test 3

```
echo "Hello from test3"
```
