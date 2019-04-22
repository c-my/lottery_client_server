from flask import Flask 
from flask_socketio import SocketIO,emit

app = Flask(__name__)

socketio = SocketIO()
socketio.init_app(app)

"""
对app进行一些路由设置
"""

"""
对socketio进行一些监听设置
"""
@socketio.on('request_for_response', namespace='/testnamespace')
def give_response(data):
    value = data.get('param')

    #进行一些对value的处理或者其他操作,在此期间可以随时会调用emit方法向前台发送消息
    emit('response',{'code':'200','msg':'start to process...'})

    time.sleep(5)
    emit('response',{'code':'200','msg':'processed'})

if __name__ == '__main__':
    socketio.run(app, debug=True, host='localhost',port=1923)
    #这里就不再用app.run而用socketio.run了。socketio.run的参数和app.run也都差不多