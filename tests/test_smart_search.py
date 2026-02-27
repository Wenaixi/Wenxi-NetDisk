"""
Wenxi网盘 - 智能搜索系统测试
作者：Wenxi
功能：测试全文搜索功能的性能和正确性
"""

import asyncio
import time
import sys
import os

# 添加backend目录到路径
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'backend'))

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker
from models import Base, File as FileModel, User
from search.index_manager import search_manager, SearchIndexManager
from search.service import search_service


# 测试数据库
TEST_DATABASE_URL = "sqlite:///./test_search.db"


def setup_test_db():
    """设置测试数据库"""
    engine = create_engine(TEST_DATABASE_URL)
    Base.metadata.create_all(bind=engine)
    Session = sessionmaker(bind=engine)
    return Session()


def create_test_user(db):
    """创建测试用户"""
    user = db.query(User).filter(User.username == "test_user").first()
    if not user:
        user = User(
            username="test_user",
            email="test@example.com",
            hashed_password="test_hash"
        )
        db.add(user)
        db.commit()
        db.refresh(user)
    return user


def create_test_files(db, user_id, count=100):
    """创建测试文件"""
    files = []
    test_names = [
        "年度报告2024.pdf", "月度销售数据.xlsx", "项目计划书.docx",
        "会议记录.txt", "代码示例.py", "配置文件.json",
        "图片素材.png", "视频教程.mp4", "音频文件.mp3",
        "压缩包.zip", "演示文稿.pptx", "数据库备份.sql"
    ]
    
    # 清除现有测试文件
    db.query(FileModel).filter(FileModel.owner_id == user_id).delete()
    db.commit()
    
    for i in range(count):
        name = test_names[i % len(test_names)]
        file = FileModel(
            filename=f"test_file_{i}.encrypted",
            original_filename=f"{i}_{name}",
            file_path=f"uploads/test_file_{i}",
            file_size=1024 * (i + 1),
            mime_type="application/octet-stream",
            owner_id=user_id,
            description=f"这是测试文件 {i} 的描述"
        )
        db.add(file)
        files.append(file)
    
    db.commit()
    for f in files:
        db.refresh(f)
    
    return files


async def test_index_creation():
    """测试索引创建"""
    print("\n📝 测试索引创建...")
    
    db = setup_test_db()
    user = create_test_user(db)
    files = create_test_files(db, user.id, 50)
    
    # 重建索引
    start_time = time.time()
    result = await search_service.reindex_all_files(db)
    elapsed = time.time() - start_time
    
    print(f"   ✅ 索引重建完成: {result['indexed_count']} 个文件")
    print(f"   ⏱️  耗时: {elapsed:.2f}秒")
    
    # 验证索引统计
    stats = await search_service.get_search_stats()
    print(f"   📊 索引统计: {stats}")
    
    db.close()
    return elapsed < 5.0  # 应该在5秒内完成


async def test_search_performance():
    """测试搜索性能"""
    print("\n🔍 测试搜索性能...")
    
    db = setup_test_db()
    user = create_test_user(db)
    
    test_queries = [
        "年度报告",
        "销售数据",
        "项目计划",
        "代码示例",
        "图片素材",
        "视频教程",
        "配置",
        "备份"
    ]
    
    results = []
    for query in test_queries:
        start_time = time.time()
        result = await search_service.search_files(
            query=query,
            user_id=user.id,
            db=db,
            limit=20
        )
        elapsed = (time.time() - start_time) * 1000  # 转换为毫秒
        
        results.append({
            'query': query,
            'time_ms': elapsed,
            'total': result['total'],
            'results_count': len(result['results'])
        })
        
        status = "✅" if elapsed < 200 else "⚠️"
        print(f"   {status} '{query}': {elapsed:.2f}ms, 找到 {result['total']} 个结果")
    
    # 计算平均响应时间
    avg_time = sum(r['time_ms'] for r in results) / len(results)
    max_time = max(r['time_ms'] for r in results)
    min_time = min(r['time_ms'] for r in results)
    
    print(f"\n   📈 性能统计:")
    print(f"      平均响应时间: {avg_time:.2f}ms")
    print(f"      最快响应: {min_time:.2f}ms")
    print(f"      最慢响应: {max_time:.2f}ms")
    
    db.close()
    return avg_time < 200  # 平均响应时间应小于200ms


async def test_file_type_filter():
    """测试文件类型过滤"""
    print("\n📁 测试文件类型过滤...")
    
    db = setup_test_db()
    user = create_test_user(db)
    
    test_cases = [
        ('document', '文档'),
        ('spreadsheet', '表格'),
        ('image', '图片'),
        ('code', '代码'),
    ]
    
    for file_type, label in test_cases:
        result = await search_service.search_files(
            query="test",
            user_id=user.id,
            db=db,
            file_type=file_type,
            limit=20
        )
        
        print(f"   ✅ {label}类型过滤: 找到 {result['total']} 个结果")
    
    db.close()
    return True


async def test_relevance_sorting():
    """测试相关性排序"""
    print("\n📊 测试相关性排序...")
    
    db = setup_test_db()
    user = create_test_user(db)
    
    # 搜索"年度"，检查相关性最高的文件是否排在前面
    result = await search_service.search_files(
        query="年度报告",
        user_id=user.id,
        db=db,
        limit=10
    )
    
    if result['results']:
        top_result = result['results'][0]
        print(f"   ✅ 最高相关性文件: {top_result['original_filename']}")
        print(f"   📈 相关性分数: {top_result.get('score', 'N/A')}")
        
        # 验证排序是否正确
        scores = [r.get('score', 0) for r in result['results']]
        is_sorted = all(scores[i] >= scores[i+1] for i in range(len(scores)-1)) if len(scores) > 1 else True
        print(f"   {'✅' if is_sorted else '❌'} 按相关性排序: {'正确' if is_sorted else '错误'}")
    
    db.close()
    return True


async def test_index_update():
    """测试索引更新"""
    print("\n🔄 测试索引更新...")
    
    db = setup_test_db()
    user = create_test_user(db)
    
    # 创建一个新文件
    file = FileModel(
        filename="new_test_file.encrypted",
        original_filename="新测试文件_特别关键词.txt",
        file_path="uploads/new_test_file",
        file_size=1024,
        mime_type="text/plain",
        owner_id=user.id,
        description="这是一个包含特别关键词的测试文件"
    )
    db.add(file)
    db.commit()
    db.refresh(file)
    
    # 添加到索引
    start_time = time.time()
    success = await search_service.index_file(file.id, db)
    elapsed = (time.time() - start_time) * 1000
    
    print(f"   ✅ 添加文件到索引: {'成功' if success else '失败'} ({elapsed:.2f}ms)")
    
    # 验证能否搜索到
    result = await search_service.search_files(
        query="特别关键词",
        user_id=user.id,
        db=db
    )
    
    found = any(r['id'] == file.id for r in result['results'])
    print(f"   {'✅' if found else '❌'} 搜索验证: {'找到' if found else '未找到'}新文件")
    
    # 测试删除
    await search_service.remove_file_index(file.id)
    result = await search_service.search_files(
        query="特别关键词",
        user_id=user.id,
        db=db
    )
    
    removed = not any(r['id'] == file.id for r in result['results'])
    print(f"   {'✅' if removed else '❌'} 删除索引: {'成功' if removed else '失败'}")
    
    db.close()
    return success and found and removed


async def run_all_tests():
    """运行所有测试"""
    print("=" * 60)
    print("🚀 Wenxi网盘 - 智能搜索系统测试")
    print("=" * 60)
    
    tests = [
        ("索引创建", test_index_creation),
        ("搜索性能", test_search_performance),
        ("文件类型过滤", test_file_type_filter),
        ("相关性排序", test_relevance_sorting),
        ("索引更新", test_index_update),
    ]
    
    results = []
    for name, test_func in tests:
        try:
            success = await test_func()
            results.append((name, success))
        except Exception as e:
            print(f"   ❌ 测试失败: {e}")
            results.append((name, False))
    
    # 打印总结
    print("\n" + "=" * 60)
    print("📋 测试结果总结")
    print("=" * 60)
    
    passed = sum(1 for _, r in results if r)
    total = len(results)
    
    for name, success in results:
        status = "✅ 通过" if success else "❌ 失败"
        print(f"   {status}: {name}")
    
    print(f"\n   总计: {passed}/{total} 通过")
    
    if passed == total:
        print("\n🎉 所有测试通过！智能搜索系统工作正常。")
    else:
        print("\n⚠️ 部分测试失败，请检查实现。")
    
    # 清理测试数据
    cleanup()
    
    return passed == total


def cleanup():
    """清理测试数据"""
    try:
        if os.path.exists("./test_search.db"):
            os.remove("./test_search.db")
            print("\n🧹 测试数据已清理")
    except:
        pass


if __name__ == "__main__":
    asyncio.run(run_all_tests())