import os
import time
import webbrowser
import win32gui
import win32con
import pystray
from pystray import Icon as icon, Menu as menu, MenuItem as item
from PIL import Image, ImageDraw

import supervisor


def create_image(width, height, color1, color2):
    # Generate an image and draw a pattern
    image = Image.new('RGB', (width, height), color1)
    dc = ImageDraw.Draw(image)
    dc.rectangle(
        (width // 2, 0, width, height // 2),
        fill=color2)
    dc.rectangle(
        (0, height // 2, width // 2, height),
        fill=color2)

    return image


state = True
program = win32gui.GetForegroundWindow()


def toggle_window_show():
    global state
    state = not state
    if state:
        win32gui.ShowWindow(program, win32con.SW_SHOW)
    else:
        win32gui.ShowWindow(program, win32con.SW_HIDE)


def open_config():
    webbrowser.open('http://127.0.0.1:5000')


def close():
    win32gui.CloseWindow(program)
    os._exit(0)


# To finally show you icon, call run
if __name__ == '__main__':
    # Update the state in `on_clicked` and return the new state in
    # a `checked` callable
    supervisor.start_supervisor()
    icon(
        'test',
        create_image(64, 64, 'black', 'white'),
        menu=menu(
            item('Show / Hide console', lambda item: toggle_window_show()),
            item('Open configure', lambda item: open_config()),
         )
         ).run()
