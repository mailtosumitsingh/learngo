client = APIClient.new("http://localhost:8080/api")
print(client:mousemove(100, 200))
print(client:click(100, 200))
print(client:movewheel(5))
print(client:sendtext(100, 200, "Hello"))
print(client:type("Hello, World!"))
print(client:findtext("Sample Text"))
print(client:gettext(10, 10, 200, 200))
print(client:screenshot())
print(client:getmousecolor(100, 200))
print(client:getmouse())

for i = 1, 10, 1 do
    print(client:mousemove(5*i, i*5))
end

