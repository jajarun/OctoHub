package io.octohub.annotation;

import java.lang.annotation.ElementType;
import java.lang.annotation.Retention;
import java.lang.annotation.RetentionPolicy;
import java.lang.annotation.Target;

/**
 * 签名验证注解
 * 标记需要进行签名验证的接口
 */
@Target({ElementType.METHOD, ElementType.TYPE})
@Retention(RetentionPolicy.RUNTIME)
public @interface SignatureAuth {
    /**
     * 是否必须签名验证，默认为true
     */
    boolean required() default true;
}
