# API Design

## Context

We need to expose our resource and actions, how to make it easy to understand and use, will be our first requirement

#### Objective

Leave the convenience to our customer, and the difficulty to ourselves

## Decision

* **资源名**：小写，并使用中划线分隔名词，如host-name
* **URI使用复数**：因为复数更适合表达整体，而单数则表示文件目录结构。
* **HTTP方法**：建议使用GET, POST, DELETE; 因为PUT逻辑包含了POST，当资源不存在的时候，创建，当存在的时候，全量更新，之中蕴含了逻辑，不像GET, 
POST简单存粹。
* **HTTP代码**：
    * 200 OK
    * 201 CREATED
    * 400 BAD REQUEST，如参数错误
    * 401 UNAUTHORIZED / 403 FORBIDDEN
    * 404 NOT FOUND 
    * 500 INTERNAL ERROR
    * 502 BAD GATEWAY 当上游系统返回的请求不合法时
    * 503 SERVICE UNAVAILABLE 当服务过载时 
    * 504 GATEWAY TIMEOUT 当上游服务超时时。建议用好HTTP代码，因为我们服务所需要的场景已经被包括，很全面。不要轻易自定义状态码，增加复杂度。
* **版本号**：面向用户提供服务，提供的资源服务不要轻易做破坏性升级，影响到服务使用者。并且一旦做为服务发布，那么服务的稳定期是因当被保证的，如facebook的API承诺至少二年有效
* **分页**：拥有默认分页设置，并且需要支持用户可自定义分页标准，为了安全和可靠性，我们需要对上限和下限进行管理
* **响应格式**：首选JSON，考虑到可读性和普遍性，以及对其它格式转换器的友好性。根目录建议用字典，方便属性追加，而不会造成破坏性更新，方便迭代管理，并用户友好
* **错误信息**：作为服务，首要考虑用户友好性，http code本身具有表意性，加上错误信息的辅助，我们因该考虑如何帮助到服务调用者，进行问题排查，及下一步动作
* **API文档说明**：可视化工具，得于用户查询和测试
* **批量数据操作**：不建议提供批量数据通过api集成，使用场景有可能是批量导入，如果该任务执行时间过长，如超过1分钟，建设使用后台任务进行处理，如果较短，则可使用rpc

## Status

Agreed

## Conclusion

* RPC
* RESTful
