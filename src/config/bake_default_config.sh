sed -n '1,/const defaultTomlConfig = `/p' config.go
cat default.toml
echo ""  # make sure there's a \n before the `
sed -n '/^`$/,$p' config.go