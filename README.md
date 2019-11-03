# Templates

 - Code must go into <% and %>
 - Expressions can appear into <%= and %>
 - Unescaped output can appear whith <%== and %>
 - Functions can be included with a header tag "<%@" at the beginning of the template

It also generates a sourcemap for errors.

For example:

```
<body>
    <%
        a := "John"
        b := a + " Doe"
    %>

    <h1>Hello <%= a %></h1>

    <h1>Hello <%== b %></h1>

</body>
```

Will generate:

```
w.write(html.encode(`<body>
    `))

        a := "John"
        b := a + " Doe"

w.write(html.encode(`

    <h1>Hello `))
w.write(a)
w.write(html.encode(`</h1>

    <h1>Hello `))
w.write(html.encode(b))
w.write(html.encode(`</h1>

</body>`))
```

Functions at the beginning of the template:

```
<%@
    func myFunc() {

    }
%>
```
