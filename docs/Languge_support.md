# Language Support for HTMLs

## 1 language support package

We shall provide support packages for different language in our portal site, together with the releases of entire lucky draw software. Users can choose whether to use these plugins or not and we shall provide a default language, like Chinese, for them. 

Users download the package they need and put it into a specific directory where our software is installed. They revise a configuration file to give the name of package. The configuration file can serve as a package controller. Here’s a proposal in JSON:
```plain
{
  "language-support": "package_name" / "none" if no package
  ...other support...
}
```
The next time when users lauch our software, the backend will detect the configuration file and try to load the package as it gives. I’ll continue describing the procedure in next chapter.

### 1.1 definations

You should first define the form of package file.

**define file surfix**
You can define whatever proper surfix as you like, such as ‘.lang’

**define variables**
You need to read `console.html` which you can get from `https://github.com/sampling-luck-draw/ConsoleFrontEnd`. Then you determine which text need to be replaced with a variable. You replace all those words with variable names you defined and give an list with pairs of variable names and their corresponding old Chinese words.

**define file content organization**
The package is like a dictionary with pairs of variable names and corresponding translations. A JSON filetype is  recommended:
```plain
{
  "language": "ENGLISH",
  "vocabulary": {
    "console_title": "sampling lucky draw",
    "start_draw_btn_text": "start draw"
    ...other pairs...
  }
}
```
And also you can define a form with chinese references:
```json
{
  "language": "ENGLISH",
  "vocabulary": {
    "console_title": {
      "translate": "sampling lucky draw",
      "reference": "三百两抽奖系统"
    },
    "start_draw_btn_text": {
      "translate": "start draw",
      "reference": "开始抽奖"
    }
    ...other pairs...
  }
}
```
After you have determined a package form, you should write an detailed description document.
### 1.2 programming
You need a program to generate language support package automatically. 

**call an API**
An Chinese word list serves as input and a word list in destinate language comes out as output. You’d better not use APIs that demands VPN which are not convenient. You can search the internet and look up in blogs. Here’s an example of Youdao APIs: `http://fanyi.youdao.com/openapi?path=data-mode`

After you get the API document, you can program with any programming language to test the interface. You can write an function with POST or GET method.

Pay attention to the generality of your program.

**reformat word list**
You need to reorganize your word list to fit the package form you defined before.

## 2 backend template rendering

Go language provide a method to render views(htmls) before they are passed to the front end so that you can pass specific variables. You can get a lot of tutorials on the internet, an example: `https://blog.csdn.net/u013210620/article/details/78522293`. For Chinese query keywords, you can use "Go语言", "Golang", "模板渲染" and "模板".

You can change the form of variables you have inserted in `console.html` to fit golang’s template syntax. And in the backend GO server file, you read the word list out of the package file (using JSON parser), and then generate the dictionary from that, and pass it to renderer finally.

remember to keep the initial Chinese word list. When user doesn’t give a package or package loading fails, use the initial list to render.

**IF YOU ARE GEEK ENOUGH**
Try not to revise initial texts in the html file. If you can implement the language support render function through just inserting some attributes in elements, it would be far more universal. If so, any existing html file can use this method to update its language support, remaining compatible with its previous version.