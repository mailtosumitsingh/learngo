client = APIClient.new("http://localhost:8080/api")


base64_image = client:screenshot()
file_path = "output.png"
save_image(base64_image, file_path)
--print(client:mousemove(100, 200))
--print(client:click(100, 200))
--print(client:movewheel(5))
print(client:sendtext(100, 200, "Hello"))
print(client:type("Hello, World!"))
print(client:findtext("Sample Text"))
print(client:gettext(10, 10, 200, 200))

print(client:getmousecolor(100, 200))
print(client:getmouse())

print(client:findimage(encode_image("output.png")))