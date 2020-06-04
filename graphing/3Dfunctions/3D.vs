#version 410 core
layout (location = 0) in vec3 aPos;
layout (location = 1) in vec3 fragColor;

//uniform vec2 aOffset;
//uniform float fragColor;
//uniform mat4 transform;

uniform mat4 model;
uniform mat4 view;
uniform mat4 projection;

out vec3 fColor;

void main()
{
    fColor = fragColor;
    gl_Position = projection * view * model * vec4(aPos, 1.0);
}