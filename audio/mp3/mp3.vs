#version 410 core
layout (location = 0) in vec2 aPos;
layout (location = 1) in vec3 fragColor;

//uniform vec2 aOffset;
//uniform float fragColor;

uniform mat4 transform;

out vec3 fColor;

void main()
{
    fColor = fragColor;
    gl_Position = transform * vec4(aPos, 0.0, 1.0);
}