/*
Copyright 2022 The Apex Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

import {
  Annotated,
  AnyType,
  BaseVisitor,
  Context,
  Kind,
  Named,
  Map,
  List,
  Optional,
} from "./deps/core.ts";
import { AliasVisitor, translations } from "./deps/typescript.ts";
import { formatComment, pascalCase } from "./deps/utils.ts";

interface ComponentDirective {
  value: string;
}

export class ComponentsVisitor extends BaseVisitor {
  public visitContextBefore(_context: Context): void {
    this
      .write(`// deno-lint-ignore-file no-explicit-any no-unused-vars ban-unused-ignore
      import { Component, DataExpr, Handler, ResourceRef, Step, ValueExpr } from "../nanobus.ts";\n\n`);
  }

  visitAlias(context: Context): void {
    super.visitAlias(context);
    const { alias } = context;
    if (
      [
        "Component",
        "ValueExpr",
        "DataExpr",
        "Handler",
        "ResourceRef",
        "Step",
      ].indexOf(alias.name) != -1
    ) {
      return;
    }
    const visitor = new AliasVisitor(this.writer);
    alias.accept(context, visitor);
    this.doClass(alias, alias);
  }

  visitEnum(context: Context): void {
    super.visitEnum(context);
    const { enum: e } = context;
    const visitor = new EnumVisitor(this.writer);
    e.accept(context, visitor);
    this.doClass(e, e);
  }

  visitTypeAfter(context: Context): void {
    const { type } = context;
    const visitor = new InterfaceVisitor(this.writer);
    type.accept(context, visitor);
    this.doClass(type, type);
  }

  doClass(named: Named, annotated: Annotated): void {
    ["initializer", "transport", "router", "middleware", "action"].forEach(
      (componentType) => {
        const a = annotated.annotation(componentType);
        if (!a) {
          return;
        }

        const name = named.name.replaceAll(
          /(Config|Configuration|Settings)$/g,
          ""
        );
        const comp = a.convert<ComponentDirective>();

        this.write(`
export function ${name}(config: ${named.name}): Component<${named.name}> {
  return {
    uses: "${comp.value}",
    with: config,
  }
}\n\n`);
      }
    );
  }
}

export class InterfaceVisitor extends BaseVisitor {
  visitTypeBefore(context: Context): void {
    super.triggerTypeBefore(context);
    const { type } = context;
    this.write(`export interface ${type.name} {\n`);
  }

  visitTypeField(context: Context): void {
    super.triggerTypeField(context);
    const { type, field } = context;
    // Name is automatically passed in.
    const comp = ["initializer", "transport", "router", "middleware", "action"].find(
      (componentType) => {
        return type.annotation(componentType) != undefined;
      });
    if (comp && field.name == "name") {
      return;
    }
    this.write(formatComment("  // ", field.description));
    let t = field.type;
    let queston = "";
    if (t.kind == Kind.Optional) {
      queston = "?";
      t = (t as Optional).type;
    }
    if (field.default) {
      queston = "?";
    }
    const et = expandType(t, true);
    this.write(`  ${field.name}${queston}: ${et};\n`);
  }

  visitTypeAfter(context: Context): void {
    this.write(`}\n\n`);
    super.triggerTypeAfter(context);
  }
}

class EnumVisitor extends BaseVisitor {
  visitEnumBefore(context: Context): void {
    super.triggerEnumsBefore(context);
    this.write(formatComment("// ", context.enum.description));
    this.write(`export enum ${context.enum.name} {\n`);
  }

  visitEnumValue(context: Context): void {
    const { enumValue } = context;
    this.write(formatComment("// ", enumValue.description));
    this.write(
      `\  ${pascalCase(enumValue.name)} = "${
        enumValue.display || enumValue.name
      }",\n`
    );
    super.triggerTypeField(context);
  }

  visitEnumAfter(context: Context): void {
    this.write(`}\n\n`);
    super.triggerEnumsAfter(context);
  }
}

/**
 * returns string of the expanded type of a node
 * @param type the type node that is being expanded
 * @param useOptional if the type that is being expanded is optional
 */
const expandType = (type: AnyType, useOptional: boolean): string => {
  switch (type.kind) {
    case Kind.Primitive:
    case Kind.Alias:
    case Kind.Enum:
    case Kind.Type:
    case Kind.Union: {
      const namedValue = (type as Named).name;
      const translation = translations.get(namedValue);
      if (translation != undefined) {
        return translation!;
      }
      if (namedValue == "Component") {
        return "Component<any>";
      }
      return namedValue;
    }
    case Kind.Map:
      return `{ [key: ${expandType((type as Map).keyType, true)}]: ${expandType(
        (type as Map).valueType,
        true
      )} }`;
    case Kind.List:
      return `${expandType((type as List).type, true)}[]`;
    case Kind.Optional: {
      const expanded = expandType((type as Optional).type, true);
      if (useOptional) {
        return `${expanded} | undefined`;
      }
      return expanded;
    }
    default:
      return "unknown";
  }
};
