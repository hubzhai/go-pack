golang ��xml�����
dom4g�ṩxml���Ĳ�����������ڵ� ���ӣ�ɾ������ѯ���������ӣ��޸ģ�ɾ������ѯ�ȹ���

�����򵥽��ܣ�

����xml�ĵ�������Elementָ��
1��LoadByStream  
2��LoadByXml   ����Ϊ�ַ���

�����ڵ�
1��LoadByStream
2��LoadByXml
3��NewElement   ����ָ��������ֵ��Elementָ��

ת�ַ������
1��ToString   ��ǰ�ڵ�xml�ַ���
2��ToXml      �����ĵ�xml�ַ���
3��SyncToXml  ��������ĵ�xml�ַ�����Ϊͬ�����������������нڵ㶼��������
4��DocLength  �����ĵ��Ľڵ��� 

��ȡ�ڵ����֣�ֵ������
1����ȡElement��Name()��Value��Attrs(���Լ���)

���Բ���
1��AttrValue  ����ָ�����ֵ����Ե�ֵ
2��AddAttr    ����ǰ�ڵ�����һ��ָ��������ֵ������
3��RemoveAttr  ɾ��ָ�����ֵ�����

�ӽڵ����
1��Node  ����ָ�����ֵ�Element�ӽڵ�
2��Nodes ����ָ�����ֵ�Element ����
3��NodesLength  �����ӽڵ����
4��AllNodes  ���������ӽڵ㼯��
5��RemoveNode ɾ��ָ�����ֵ��ӽڵ�(�����ж����ͬ���ֵĽڵ㣬������ɾ��)
6��AddNode  ����һ���ӽڵ�
7��AddNodeByString  ����һ���ӽڵ㣬����Ϊ�ַ����磺<a>b</a>  �ṹ��Ϊxml�ṹ

��ȡ���ڵ�
1��Parent  ���ظ��ڵ�Elementָ�룬����ǰ�ڵ�Ϊ���ڵ㣬�򷵻�nil
 

����������Լ������ļ�dom_test.go